package taskvault

import (
	"crypto/tls"
	"errors"
	"expvar"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/armon/go-metrics"
	"github.com/danluki/taskvault/pkg/types"
	"github.com/hashicorp/memberlist"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"github.com/hashicorp/serf/serf"
	"github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
	"golang.org/x/exp/rand"
)

const (
	raftTimeout      = 30 * time.Second
	raftLogCacheSize = 512
	minRaftProtocol  = 3
)

var (
	expNode             = expvar.NewString("node")
	ErrLeaderNotFound   = errors.New("no member leader found")
	ErrNoSuitableServer = errors.New("no suitable server found")
)

type Node = serf.Member

type Agent struct {
	Store      Storage
	TLSConfig  *tls.Config
	config     *Config
	eventCh    chan serf.Event
	shutdownCh chan struct{}
	ready      bool

	raftTransport *raft.NetworkTransport
	raft          *raft.Raft
	serf          *serf.Serf
	HTTPTransport Transport
	raftStore     RaftStore
	GRPCClient    TaskvaultGRPCClient
	raftLayer     *RaftLayer
	reconcileCh   chan serf.Member
	raftInmem     *raft.InmemStore
	GRPCServer    TaskvaultGRPCServer
	peers         map[string][]*ServerParts
	logger        *logrus.Entry
	retryJoinCh   chan error
	leaderCh      <-chan bool
	localPeers    map[raft.ServerAddress]*ServerParts
	peerLock      sync.RWMutex
	serverLookup  *ServerLookup
	listener      net.Listener
}

type RaftStore interface {
	raft.StableStore
	raft.LogStore
	Close() error
}

type AgentOption func(agent *Agent)

func NewAgent(config *Config, options ...AgentOption) *Agent {
	agent := &Agent{
		config:       config,
		retryJoinCh:  make(chan error),
		serverLookup: NewServerLookup(),
	}

	for _, option := range options {
		option(agent)
	}

	return agent
}

func (a *Agent) Start() error {
	log := InitLogger(a.config.LogLevel, a.config.NodeName)
	a.logger = log

	if err := a.config.normalizeAddrs(); err != nil && !errors.Is(
		err, ErrResolvingHost,
	) {
		return err
	}

	s, err := a.setupSerf()
	if err != nil {
		return fmt.Errorf("agent: Can not setup serf, %s", err)
	}
	a.serf = s

	if len(a.config.RetryJoinLAN) > 0 {
		a.retryJoinLAN()
	} else {
		_, err := a.join(a.config.StartJoin, true)
		if err != nil {
			a.logger.WithError(err).WithField(
				"servers", a.config.StartJoin,
			).Warn("agent: Can not join")
		}
	}

	if err := initMetrics(a); err != nil {
		a.logger.Fatal("agent: Can not setup metrics")
	}

	expNode.Set(a.config.NodeName)

	if a.config.AdvertiseRPCPort <= 0 {
		a.config.AdvertiseRPCPort = a.config.RPCPort
	}

	addr := a.bindRPCAddr()
	l, err := net.Listen("tcp", addr)
	if err != nil {
		a.logger.Fatal(err)
	}
	a.listener = l

	a.StartServer()

	if a.GRPCClient == nil {
		a.GRPCClient = NewGRPCClient(nil, a, a.logger)
	}

	tags := a.serf.LocalMember().Tags
	tags["rpc_addr"] = a.advertiseRPCAddr()
	tags["port"] = strconv.Itoa(a.config.AdvertiseRPCPort)
	if err := a.serf.SetTags(tags); err != nil {
		return fmt.Errorf("agent: Error setting tags: %w", err)
	}

	go a.eventLoop()
	a.ready = true

	return nil
}

func (a *Agent) RetryJoinCh() <-chan error {
	return a.retryJoinCh
}

func (a *Agent) JoinLAN(addrs []string) (int, error) {
	return a.serf.Join(addrs, true)
}

func (a *Agent) Stop() error {
	a.logger.Info("agent: Called member stop, now stopping")

	// TODO: Check why Shutdown().Error() is not working
	_ = a.raft.Shutdown()

	if err := a.Store.Shutdown(); err != nil {
		return err
	}

	if err := a.serf.Leave(); err != nil {
		return err
	}

	if err := a.serf.Shutdown(); err != nil {
		return err
	}

	return nil
}

// UpdateTags updates the tag configuration for this agent
func (a *Agent) UpdateTags(tags map[string]string) {
	currentTags := a.serf.LocalMember().Tags
	for _, tagName := range []string{
		"role", "version", "server", "bootstrap", "expect", "port", "rpc_addr",
	} {
		if val, exists := currentTags[tagName]; exists {
			tags[tagName] = val
		}
	}
	tags["dc"] = a.config.Datacenter
	tags["region"] = a.config.Region

	err := a.serf.SetTags(tags)
	if err != nil {
		a.logger.Warnf("Setting tags unsuccessful: %s.", err.Error())
	}
}

func (a *Agent) setupRaft() error {
	if a.config.BootstrapExpect > 0 {
		if a.config.BootstrapExpect == 1 {
			a.config.Bootstrap = true
		}
	}

	logger := ioutil.Discard
	if a.logger.Logger.Level == logrus.DebugLevel {
		logger = a.logger.Logger.Writer()
	}

	transConfig := &raft.NetworkTransportConfig{
		Stream:                a.raftLayer,
		MaxPool:               3,
		Timeout:               raftTimeout,
		ServerAddressProvider: a.serverLookup,
	}
	transport := raft.NewNetworkTransportWithConfig(transConfig)
	a.raftTransport = transport

	config := raft.DefaultConfig()

	// Raft performance
	raftMultiplier := a.config.RaftMultiplier
	if raftMultiplier < 1 || raftMultiplier > 10 {
		return fmt.Errorf(
			"raft-multiplier cannot be %d. Must be between 1 and 10",
			raftMultiplier,
		)
	}
	config.HeartbeatTimeout = config.HeartbeatTimeout * time.Duration(raftMultiplier)
	config.ElectionTimeout = config.ElectionTimeout * time.Duration(raftMultiplier)
	config.LeaderLeaseTimeout = config.LeaderLeaseTimeout * time.Duration(a.config.RaftMultiplier)

	config.LogOutput = logger
	config.LocalID = raft.ServerID(a.config.NodeName)

	var logStore raft.LogStore
	var stableStore raft.StableStore
	var snapshots raft.SnapshotStore
	if a.config.DevMode {
		store := raft.NewInmemStore()
		a.raftInmem = store
		stableStore = store
		logStore = store
		snapshots = raft.NewDiscardSnapshotStore()
	} else {
		var err error

		snapshots, err = raft.NewFileSnapshotStore(
			filepath.Join(
				a.config.DataDir, "raft",
			), 3, logger,
		)
		if err != nil {
			return fmt.Errorf("file snapshot store: %s", err)
		}

		if a.raftStore == nil {
			s, err := raftboltdb.NewBoltStore(
				filepath.Join(
					a.config.DataDir, "raft", "raft.db",
				),
			)
			if err != nil {
				return fmt.Errorf("error creating new raft store: %s", err)
			}
			a.raftStore = s
		}
		stableStore = a.raftStore

		cacheStore, err := raft.NewLogCache(raftLogCacheSize, a.raftStore)
		if err != nil {
			a.raftStore.Close()
			return err
		}
		logStore = cacheStore

		// Check for peers.json file for recovery
		peersFile := filepath.Join(a.config.DataDir, "raft", "peers.json")
		if _, err := os.Stat(peersFile); err == nil {
			a.logger.Info("found peers.json file, recovering Raft configuration...")
			var configuration raft.Configuration
			configuration, err = raft.ReadConfigJSON(peersFile)
			if err != nil {
				return fmt.Errorf(
					"recovery failed to parse peers.json: %v", err,
				)
			}
			store, err := NewStore(a.logger)
			if err != nil {
				a.logger.WithError(err).Fatal("taskvault: Error initializing store")
			}
			tmpFsm := newFSM(store, a.logger)
			if err := raft.RecoverCluster(
				config, tmpFsm,
				logStore, stableStore, snapshots, transport, configuration,
			); err != nil {
				return fmt.Errorf("recovery failed: %v", err)
			}
			if err := os.Remove(peersFile); err != nil {
				return fmt.Errorf(
					"recovery failed to delete peers.json, please delete manually (see peers.info for details): %v",
					err,
				)
			}
			a.logger.Info("deleted peers.json file after successful recovery")
		}
	}

	if a.config.Bootstrap || a.config.DevMode {
		hasState, err := raft.HasExistingState(logStore, stableStore, snapshots)
		if err != nil {
			return err
		}
		if !hasState {
			configuration := raft.Configuration{
				Servers: []raft.Server{
					{
						ID:      config.LocalID,
						Address: transport.LocalAddr(),
					},
				},
			}
			if err := raft.BootstrapCluster(
				config, logStore, stableStore, snapshots, transport,
				configuration,
			); err != nil {
				return err
			}
		}
	}

	fsm := newFSM(a.Store, a.logger)
	rft, err := raft.NewRaft(
		config, fsm, logStore, stableStore, snapshots, transport,
	)
	if err != nil {
		return fmt.Errorf("new raft: %s", err)
	}
	a.leaderCh = rft.LeaderCh()
	a.raft = rft

	return nil
}

// setupSerf is used to create the agent we use
func (a *Agent) setupSerf() (*serf.Serf, error) {
	config := a.config

	a.localPeers = make(map[raft.ServerAddress]*ServerParts)
	a.peers = make(map[string][]*ServerParts)

	bindIP, bindPort, err := config.AddrParts(config.BindAddr)
	if err != nil {
		return nil, fmt.Errorf("invalid bind address: %s", err)
	}

	var advertiseIP string
	var advertisePort int
	if config.AdvertiseAddr != "" {
		advertiseIP, advertisePort, err = config.AddrParts(config.AdvertiseAddr)
		if err != nil {
			return nil, fmt.Errorf("invalid advertise address: %s", err)
		}
	}

	encryptKey, err := config.EncryptBytes()
	if err != nil {
		return nil, fmt.Errorf("invalid encryption key: %s", err)
	}

	serfConfig := serf.DefaultConfig()
	serfConfig.Init()

	serfConfig.Tags = a.config.Tags
	serfConfig.Tags["role"] = "taskvault"
	serfConfig.Tags["dc"] = a.config.Datacenter
	serfConfig.Tags["region"] = a.config.Region
	serfConfig.Tags["version"] = Version
	serfConfig.Tags["server"] = strconv.FormatBool(true)
	if a.config.Bootstrap {
		serfConfig.Tags["bootstrap"] = "1"
	}
	if a.config.BootstrapExpect != 0 {
		serfConfig.Tags["expect"] = fmt.Sprintf("%d", a.config.BootstrapExpect)
	}

	switch config.Profile {
	case "lan":
		serfConfig.MemberlistConfig = memberlist.DefaultLANConfig()
	case "wan":
		serfConfig.MemberlistConfig = memberlist.DefaultWANConfig()
	case "local":
		serfConfig.MemberlistConfig = memberlist.DefaultLocalConfig()
	default:
		return nil, fmt.Errorf("unknown profile: %s", config.Profile)
	}

	serfConfig.MemberlistConfig.BindAddr = bindIP
	serfConfig.MemberlistConfig.BindPort = bindPort
	serfConfig.MemberlistConfig.AdvertiseAddr = advertiseIP
	serfConfig.MemberlistConfig.AdvertisePort = advertisePort
	serfConfig.MemberlistConfig.SecretKey = encryptKey
	serfConfig.NodeName = config.NodeName
	serfConfig.Tags = config.Tags
	serfConfig.CoalescePeriod = 3 * time.Second
	serfConfig.QuiescentPeriod = time.Second
	serfConfig.UserCoalescePeriod = 3 * time.Second
	serfConfig.UserQuiescentPeriod = time.Second
	serfConfig.ReconnectTimeout, err = time.ParseDuration(config.SerfReconnectTimeout)

	if err != nil {
		a.logger.Fatal(err)
	}

	a.eventCh = make(chan serf.Event, 2048)
	serfConfig.EventCh = a.eventCh

	a.logger.Info("agent: taskvault agent starting")

	if a.logger.Logger.Level == logrus.DebugLevel {
		serfConfig.LogOutput = a.logger.Logger.Writer()
		serfConfig.MemberlistConfig.LogOutput = a.logger.Logger.Writer()
	} else {
		serfConfig.LogOutput = ioutil.Discard
		serfConfig.MemberlistConfig.LogOutput = ioutil.Discard
	}

	// Create serf first
	serf, err := serf.Create(serfConfig)
	if err != nil {
		a.logger.Error(err)
		return nil, err
	}
	return serf, nil
}

func (a *Agent) Config() *Config {
	return a.config
}

func (a *Agent) SetConfig(c *Config) {
	a.config = c
}

func (a *Agent) StartServer() {
	if a.Store == nil {
		s, err := NewStore(a.logger)
		if err != nil {
			a.logger.WithError(err).Fatal("taskvault: Error initializing store")
		}
		a.Store = s
	}

	if a.HTTPTransport == nil {
		a.HTTPTransport = NewTransport(a, a.logger)
	}
	a.HTTPTransport.ServeHTTP()

	tcpm := cmux.New(a.listener)
	var grpcl, raftl net.Listener

	if a.TLSConfig != nil {
		a.raftLayer = NewTLSRaftLayer(a.TLSConfig, a.logger)

		tlsl := tcpm.Match(cmux.Any())
		tlsl = tls.NewListener(tlsl, a.TLSConfig)

		tlsm := cmux.New(tlsl)

		grpcl = tlsm.MatchWithWriters(
			cmux.HTTP2MatchHeaderFieldSendSettings(
				"content-type", "application/grpc",
			),
		)

		raftl = tlsm.Match(cmux.Any())

		go func() {
			if err := tlsm.Serve(); err != nil {
				a.logger.Fatal(err)
			}
		}()
	} else {
		a.raftLayer = NewRaftLayer(a.logger)

		grpcl = tcpm.MatchWithWriters(
			cmux.HTTP2MatchHeaderFieldSendSettings(
				"content-type", "application/grpc",
			),
		)

		raftl = tcpm.Match(cmux.Any())
	}

	if a.GRPCServer == nil {
		a.GRPCServer = NewGRPCServer(a, a.logger)
	}

	if err := a.GRPCServer.Serve(grpcl); err != nil {
		a.logger.WithError(err).Fatal("agent: RPC server failed to start")
	}

	if err := a.raftLayer.Open(raftl); err != nil {
		a.logger.Fatal(err)
	}

	if err := a.setupRaft(); err != nil {
		a.logger.WithError(err).Fatal("agent: Raft layer failed to start")
	}

	// Start serving everything
	go func() {
		if err := tcpm.Serve(); err != nil {
			a.logger.Fatal(err)
		}
	}()
	go a.monitorLeadership()
}

func (a *Agent) leaderMember() (*serf.Member, error) {
	l := a.raft.Leader()
	for _, member := range a.serf.Members() {
		if member.Tags["rpc_addr"] == string(l) {
			return &member, nil
		}
	}
	return nil, ErrLeaderNotFound
}

func (a *Agent) IsLeader() bool {
	return a.raft.State() == raft.Leader
}

func (a *Agent) Members() []serf.Member {
	return a.serf.Members()
}

func (a *Agent) LocalMember() serf.Member {
	return a.serf.LocalMember()
}

func (a *Agent) Leader() raft.ServerAddress {
	return a.raft.Leader()
}

func (a *Agent) Servers() (members []*ServerParts) {
	for _, member := range a.serf.Members() {
		ok, parts := isServer(member)
		if !ok || member.Status != serf.StatusAlive {
			continue
		}
		members = append(members, parts)
	}
	return members
}

func (a *Agent) LocalServers() (members []*ServerParts) {
	for _, member := range a.serf.Members() {
		ok, parts := isServer(member)
		if !ok || member.Status != serf.StatusAlive {
			continue
		}
		if a.config.Region == parts.Region {
			members = append(members, parts)
		}
	}
	return members
}

func (a *Agent) eventLoop() {
	serfShutdownCh := a.serf.ShutdownCh()
	a.logger.Info("agent: Listen for events")
	for {
		select {
		case e := <-a.eventCh:
			a.logger.WithField(
				"event", e.String(),
			).Info("agent: Received event")
			metrics.IncrCounter(
				[]string{"agent", "event_received", e.String()}, 1,
			)

			if me, ok := e.(serf.MemberEvent); ok {
				for _, member := range me.Members {
					a.logger.WithFields(
						logrus.Fields{
							"node":   a.config.NodeName,
							"member": member.Name,
							"event":  e.EventType(),
						},
					).Debug("agent: Member event")
				}

				// serfEventHandler is used to handle events from the serf cluster
				switch e.EventType() {
				case serf.EventMemberJoin:
					a.nodeJoin(me)
					a.localMemberEvent(me)
				case serf.EventMemberLeave, serf.EventMemberFailed:
					a.nodeFailed(me)
					a.localMemberEvent(me)
				case serf.EventMemberReap:
					a.localMemberEvent(me)
				case serf.EventMemberUpdate:
					a.lanNodeUpdate(me)
					a.localMemberEvent(me)
				case serf.EventUser, serf.EventQuery: // Ignore
				default:
					a.logger.WithField(
						"event", e.String(),
					).Warn("agent: Unhandled serf event")
				}
			}

		case <-serfShutdownCh:
			a.logger.Warn("agent: Serf shutdown detected, quitting")
			return
		}
	}
}

func (a *Agent) join(addrs []string, replay bool) (n int, err error) {
	a.logger.Infof("agent: joining: %v replay: %v", addrs, replay)
	n, err = a.serf.Join(addrs, !replay)
	if n > 0 {
		a.logger.Infof("agent: joined: %d nodes", n)
	}
	if err != nil {
		a.logger.Warnf("agent: error joining: %v", err)
	}
	return
}

func (a *Agent) getTargetNodes(
	tags map[string]string, selectFunc func([]Node) int,
) []Node {
	bareTags, cardinality := cleanTags(tags, a.logger)
	nodes := a.getQualifyingNodes(a.serf.Members(), bareTags)
	return selectNodes(nodes, cardinality, selectFunc)
}

func (a *Agent) getQualifyingNodes(
	nodes []Node, bareTags map[string]string,
) []Node {
	// Determine the usable set of nodes
	qualifiers := filterArray(
		nodes, func(node Node) bool {
			return node.Status == serf.StatusAlive &&
				node.Tags["region"] == a.config.Region &&
				nodeMatchesTags(node, bareTags)
		},
	)
	return qualifiers
}

func defaultSelector(nodes []Node) int {
	return rand.Intn(len(nodes))
}

func selectNodes(
	nodes []Node, cardinality int, selectFunc func([]Node) int,
) []Node {
	numNodes := len(nodes)
	if numNodes <= cardinality {
		return nodes
	}

	for ; cardinality > 0; cardinality-- {
		// Select a node
		chosenIndex := selectFunc(nodes[:numNodes])

		nodes[numNodes-1], nodes[chosenIndex] = nodes[chosenIndex], nodes[numNodes-1]
		numNodes--
	}

	return nodes[numNodes:]
}

func filterArray(arr []Node, filterFunc func(Node) bool) []Node {
	for i := len(arr) - 1; i >= 0; i-- {
		if !filterFunc(arr[i]) {
			arr[i] = arr[len(arr)-1]
			arr = arr[:len(arr)-1]
		}
	}
	return arr
}

func (a *Agent) advertiseRPCAddr() string {
	bindIP := a.serf.LocalMember().Addr
	return net.JoinHostPort(
		bindIP.String(), strconv.Itoa(a.config.AdvertiseRPCPort),
	)
}

func (a *Agent) bindRPCAddr() string {
	bindIP, _, _ := a.config.AddrParts(a.config.BindAddr)
	return net.JoinHostPort(bindIP, strconv.Itoa(a.config.RPCPort))
}

func (a *Agent) applySetPair(pair *types.Pair) error {
	cmd, err := Encode(AddPairType, pair)
	if err != nil {
		return err
	}
	af := a.raft.Apply(cmd, raftTimeout)
	if err := af.Error(); err != nil {
		return err
	}

	return nil
}

// RaftApply applies a command to the Raft log
func (a *Agent) RaftApply(cmd []byte) raft.ApplyFuture {
	return a.raft.Apply(cmd, raftTimeout)
}

// Check if the server is alive and select it
func (a *Agent) checkAndSelectServer() (string, error) {
	var peers []string
	for _, p := range a.LocalServers() {
		peers = append(peers, p.RPCAddr.String())
	}

	for _, peer := range peers {
		a.logger.WithField("peer", peer).Debug("Checking peer")
		conn, err := net.DialTimeout("tcp", peer, 1*time.Second)
		if err == nil {
			conn.Close()
			a.logger.WithField("peer", peer).Debug("Found good peer")
			return peer, nil
		}
	}
	return "", ErrNoSuitableServer
}
