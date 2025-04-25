package taskvault

import (
    "fmt"
    "net"
    "sync"
    "time"

    metrics "github.com/hashicorp/go-metrics"
    "github.com/hashicorp/raft"
    "github.com/hashicorp/serf/serf"
)

const (
    barrierWriteTimeout = 2 * time.Minute
)

func (a *Agent) monitorLeadership() {
    var weAreLeaderCh chan struct{}
    var leaderLoop sync.WaitGroup
    for {
        a.logger.Info("taskvault: monitoring leadership")
        select {
        case isLeader := <-a.leaderCh:
            switch {
            case isLeader:
                if weAreLeaderCh != nil {
                    a.logger.Error("taskvault: attempted to start the leader loop while running")
                    continue
                }

                weAreLeaderCh = make(chan struct{})
                leaderLoop.Add(1)
                go func(ch chan struct{}) {
                    defer leaderLoop.Done()
                    a.leaderLoop(ch)
                }(weAreLeaderCh)
                a.logger.Info("taskvault: cluster leadership acquired")

            default:
                if weAreLeaderCh == nil {
                    a.logger.Error("taskvault: attempted to stop the leader loop while not running")
                    continue
                }

                a.logger.Debug("taskvault: shutting down leader loop")
                close(weAreLeaderCh)
                leaderLoop.Wait()
                weAreLeaderCh = nil
                a.logger.Info("taskvault: cluster leadership lost")
            }

        case <-a.shutdownCh:
            return
        }
    }
}

func (a *Agent) leaderLoop(stopCh chan struct{}) {
    var reconcileCh chan serf.Member

RECONCILE:
    reconcileCh = nil
    interval := time.After(a.config.ReconcileInterval)

    start := time.Now()
    barrier := a.raft.Barrier(barrierWriteTimeout)
    if err := barrier.Error(); err != nil {
        a.logger.WithError(err).Error("taskvault: failed to wait for barrier")
        goto WAIT
    }
    metrics.MeasureSince([]string{"taskvault", "leader", "barrier"}, start)

    if err := a.reconcile(); err != nil {
        a.logger.WithError(err).Error("taskvault: failed to reconcile")
        goto WAIT
    }

    reconcileCh = a.reconcileCh

    select {
    case <-stopCh:
        return
    default:
    }

WAIT:
    for {
        select {
        case <-stopCh:
            return
        case <-a.shutdownCh:
            return
        case <-interval:
            goto RECONCILE
        case member := <-reconcileCh:
            if err := a.reconcileMember(member); err != nil {
                a.logger.WithError(err).Error("taskvault: failed to reconcile member")
            }
        }
    }
}

func (a *Agent) reconcile() error {
    defer metrics.MeasureSince(
        []string{"taskvault", "leader", "reconcile"}, time.Now(),
    )

    members := a.serf.Members()
    for _, member := range members {
        if err := a.reconcileMember(member); err != nil {
            return err
        }
    }
    return nil
}

func (a *Agent) reconcileMember(member serf.Member) error {
    valid, parts := isServer(member)
    if !valid {
        return nil
    }
    defer metrics.MeasureSince(
        []string{
            "taskvault", "leader", "reconcileMember",
        }, time.Now(),
    )

    var err error
    switch member.Status {
    case serf.StatusAlive:
        err = a.addRaftPeer(member, parts)
    case serf.StatusLeft:
        err = a.removeRaftPeer(member, parts)
    }
    if err != nil {
        a.logger.WithError(err).WithField(
            "member", member,
        ).Error("failed to reconcile member")
        return err
    }
    return nil
}

func (a *Agent) addRaftPeer(m serf.Member, parts *ServerParts) error {
    members := a.serf.Members()
    if parts.Bootstrap {
        for _, member := range members {
            valid, p := isServer(member)
            if valid && member.Name != m.Name && p.Bootstrap {
                a.logger.Errorf(
                    "taskvault: '%v' and '%v' are both in bootstrap mode. Only one node should be in bootstrap mode, not adding Raft peer.",
                    m.Name,
                    member.Name,
                )
                return nil
            }
        }
    }

    addr := (&net.TCPAddr{IP: m.Addr, Port: parts.Port}).String()
    configFuture := a.raft.GetConfiguration()
    if err := configFuture.Error(); err != nil {
        a.logger.WithError(err).Error("taskvault: failed to get raft configuration")
        return err
    }

    if m.Name == a.config.NodeName {
        if l := len(configFuture.Configuration().Servers); l < 3 {
            a.logger.WithField("peer", m.Name).
                Debug("taskvault: Skipping self join check since the cluster is too small")
            return nil
        }
    }

    for _, server := range configFuture.Configuration().Servers {
        if server.Address == raft.ServerAddress(addr) || server.ID == raft.ServerID(parts.ID) {
            if server.Address == raft.ServerAddress(addr) && server.ID == raft.ServerID(parts.ID) {
                return nil
            }
            if server.Address == raft.ServerAddress(addr) {
                future := a.raft.RemoveServer(server.ID, 0, 0)
                if err := future.Error(); err != nil {
                    return fmt.Errorf(
                        "error removing server with duplicate address %q: %s",
                        server.Address,
                        err,
                    )
                }
                a.logger.WithField("server", server.Address).
                    Info("taskvault: removed server with duplicate address")
            }
        }
    }

    switch {
    case minRaftProtocol >= 3:
        addFuture := a.raft.AddVoter(
            raft.ServerID(parts.ID), raft.ServerAddress(addr), 0, 0,
        )
        if err := addFuture.Error(); err != nil {
            a.logger.WithError(err).Error("taskvault: failed to add raft peer")
            return err
        }
    }

    return nil
}

func (a *Agent) removeRaftPeer(m serf.Member, parts *ServerParts) error {
    if m.Name == a.config.NodeName {
        a.logger.Warn(
            "removing self should be done by follower", "name",
            a.config.NodeName,
        )
        return nil
    }

    configFuture := a.raft.GetConfiguration()
    if err := configFuture.Error(); err != nil {
        a.logger.WithError(err).Error("taskvault: failed to get raft configuration")
        return err
    }

    for _, server := range configFuture.Configuration().Servers {
        if server.ID == raft.ServerID(parts.ID) {
            a.logger.WithField(
                "server", server.ID,
            ).Info("taskvault: removing server by ID")
            future := a.raft.RemoveServer(raft.ServerID(parts.ID), 0, 0)
            if err := future.Error(); err != nil {
                a.logger.WithError(err).
                    WithField("server", server.ID).
                    Error("taskvault: failed to remove raft peer")
                return err
            }
            break
        }
    }

    return nil
}
