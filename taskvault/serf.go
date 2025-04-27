package taskvault

import (
	"strings"
	"time"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/serf/serf"
	"go.uber.org/zap"
)

const (
	StatusReap = serf.MemberStatus(-1)

	maxPeerRetries = 6
)

func (a *Agent) nodeJoin(me serf.MemberEvent, checkBootstrap bool) {
	for _, m := range me.Members {
		parts := toServerPart(m)
		if parts == nil {
			continue
		}

		a.serverLookup.AddServer(parts)

		if checkBootstrap {
			if a.config.BootstrapExpect != 0 {
				a.maybeBootstrap()
			}
		}
	}
}

func (a *Agent) maybeBootstrap() {
	var index uint64
	var err error
	if a.raftStore != nil {
		index, err = a.raftStore.LastIndex()
	} else {
		panic("raftStore is uninitialized")
	}
	if err != nil {
		return
	}

	if index != 0 {
		a.config.BootstrapExpect = 0
		return
	}

	members := a.serf.Members()
	var servers []ServerParts
	voters := 0
	for _, member := range members {
		parts := toServerPart(member)
		if parts == nil {
			continue
		}
		if parts.Expect != 0 && parts.Expect != a.config.BootstrapExpect {
			return
		}
		if parts.Bootstrap {
			return
		}
		voters++
		servers = append(servers, *parts)
	}

	if voters < a.config.BootstrapExpect {
		return
	}

	for _, server := range servers {
		var peers []string

		for i := range maxPeerRetries {
			configuration, err := a.GRPCClient.RaftGetConfiguration(server.RPCAddr.String())
			if err != nil {
				next := (1 << i) * time.Second
				a.logger.Error(
					"Failed retry.",
					"server", server.Name,
					"retry_interval", next.String(),
					"error", err,
				)
				time.Sleep(next)
			} else {
				for _, peer := range configuration.Servers {
					peers = append(peers, peer.Id)
				}
				break
			}
		}

		if len(peers) > 0 {
			a.logger.Info(
				"Disabling bootstrap mode, cus of existing Raft cluster",
				"server",
				server.Name,
			)
			a.config.BootstrapExpect = 0
			return
		}
	}

	var configuration raft.Configuration
	var addrs []string

	for _, server := range servers {
		addr := server.Addr.String()
		addrs = append(addrs, addr)
		id := raft.ServerID(server.ID)
		suffrage := raft.Voter
		peer := raft.Server{
			ID:       id,
			Address:  raft.ServerAddress(addr),
			Suffrage: suffrage,
		}
		configuration.Servers = append(configuration.Servers, peer)
	}
	a.logger.Info(
		"agent: attempting to bootstrap cluster...",
		"peers", strings.Join(addrs, ","),
	)
	future := a.raft.BootstrapCluster(configuration)
	if err := future.Error(); err != nil {
		a.logger.Error("agent: failed bootstrap", zap.Error(err))
	}

	a.config.BootstrapExpect = 0
}

func (a *Agent) nodeFailed(me serf.MemberEvent) {
	for _, m := range me.Members {
		parts := toServerPart(m)
		if parts == nil {
			continue
		}
		a.logger.Info("removing server ", parts)

		a.serverLookup.RemoveServer(parts)
	}
}

func (a *Agent) reapEvent(me serf.MemberEvent) {
	if !a.IsLeader() {
		return
	}

	for _, m := range me.Members {
		if me.EventType() == serf.EventMemberReap {
			m.Status = StatusReap
		}
		select {
		case a.refreshCh <- m:
		default:
		}
	}
}
