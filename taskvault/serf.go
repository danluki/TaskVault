package taskvault

import (
	"strings"
	"time"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/serf/serf"
)

const (
	StatusReap = serf.MemberStatus(-1)

	maxPeerRetries = 6
)

func (a *Agent) nodeJoin(me serf.MemberEvent) {
	for _, m := range me.Members {
		ok, parts := isServer(m)
		if !ok {
			a.logger.WithField(
				"member", m.Name,
			).Warn("non-server in gossip pool")
			continue
		}
		a.logger.WithField(
			"server", parts.Name,
		).Info("Adding LAN adding server")
		a.serverLookup.AddServer(parts)
		found := false
		a.peerLock.Lock()
		existing := a.peers[parts.Region]
		for idx, e := range existing {
			if e.Name == parts.Name {
				existing[idx] = parts
				found = true
				break
			}
		}

		if !found {
			a.peers[parts.Region] = append(existing, parts)
		}

		if parts.Region == a.config.Region {
			a.localPeers[raft.ServerAddress(parts.Addr.String())] = parts
		}
		a.peerLock.Unlock()

		if a.config.BootstrapExpect != 0 {
			a.maybeBootstrap()
		}
	}
}

func (a *Agent) maybeBootstrap() {
	var index uint64
	var err error
	if a.raftStore != nil {
		index, err = a.raftStore.LastIndex()
	} else if a.raftInmem != nil {
		index, err = a.raftInmem.LastIndex()
	} else {
		panic("neither raftInmem or raftStore is initialized")
	}
	if err != nil {
		a.logger.WithError(err).Error("failed to read last raft index")
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
		valid, p := isServer(member)
		if !valid {
			continue
		}
		if p.Region != a.config.Region {
			continue
		}
		if p.Expect != 0 && p.Expect != a.config.BootstrapExpect {
			a.logger.WithField("member", member).
				Error("peer has a conflicting expect value. All nodes should expect the same number")
			return
		}
		if p.Bootstrap {
			a.logger.WithField(
				"member", member,
			).Error("peer has bootstrap mode. Expect disabled")
			return
		}
		if valid {
			voters++
		}
		servers = append(servers, *p)
	}

	if voters < a.config.BootstrapExpect {
		return
	}

	for _, server := range servers {
		var peers []string

		for attempt := uint(0); attempt < maxPeerRetries; attempt++ {
			configuration, err := a.GRPCClient.RaftGetConfiguration(server.RPCAddr.String())
			if err != nil {
				nextRetry := (1 << attempt) * time.Second
				a.logger.Error(
					"Failed to confirm peer status for server (will retry).",
					"server", server.Name,
					"retry_interval", nextRetry.String(),
					"error", err,
				)
				time.Sleep(nextRetry)
			} else {
				for _, peer := range configuration.Servers {
					peers = append(peers, peer.Id)
				}
				break
			}
		}

		if len(peers) > 0 {
			a.logger.Info(
				"Existing Raft peers reported by server, disabling bootstrap mode",
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
		"agent: found expected number of peers, attempting to bootstrap cluster...",
		"peers", strings.Join(addrs, ","),
	)
	future := a.raft.BootstrapCluster(configuration)
	if err := future.Error(); err != nil {
		a.logger.WithError(err).Error("agent: failed to bootstrap cluster")
	}

	a.config.BootstrapExpect = 0
}

func (a *Agent) nodeFailed(me serf.MemberEvent) {
	for _, m := range me.Members {
		ok, parts := isServer(m)
		if !ok {
			continue
		}
		a.logger.Info("removing server ", parts)

		a.peerLock.Lock()
		existing := a.peers[parts.Region]
		n := len(existing)
		for i := 0; i < n; i++ {
			if existing[i].Name == parts.Name {
				existing[i], existing[n-1] = existing[n-1], nil
				existing = existing[:n-1]
				n--
				break
			}
		}

		if n == 0 {
			delete(a.peers, parts.Region)
		} else {
			a.peers[parts.Region] = existing
		}

		if parts.Region == a.config.Region {
			delete(a.localPeers, raft.ServerAddress(parts.Addr.String()))
		}
		a.peerLock.Unlock()
		a.serverLookup.RemoveServer(parts)
	}
}

func (a *Agent) localMemberEvent(me serf.MemberEvent) {
	if !a.IsLeader() {
		return
	}

	isReap := me.EventType() == serf.EventMemberReap

	for _, m := range me.Members {
		if isReap {
			m.Status = StatusReap
		}
		select {
		case a.reconcileCh <- m:
		default:
		}
	}
}

func (a *Agent) lanNodeUpdate(me serf.MemberEvent) {
	for _, m := range me.Members {
		ok, parts := isServer(m)
		if !ok {
			continue
		}
		a.logger.WithField("server", parts.String()).Info("Updating LAN server")


		a.serverLookup.AddServer(parts)
	}
}
