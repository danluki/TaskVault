package taskvault

import (
    "bytes"
    "context"
    "fmt"
    "net"
    "time"

    "github.com/armon/go-metrics"
    types2 "github.com/danluki/taskvault/pkg/types"
    "github.com/hashicorp/raft"
    "github.com/hashicorp/serf/serf"
    "github.com/sirupsen/logrus"
    "google.golang.org/grpc"
    "google.golang.org/protobuf/proto"
    "google.golang.org/protobuf/types/known/emptypb"
)

type TaskvaultGRPCServer interface {
    types2.TaskvaultServer
    Serve(net.Listener) error
}

type GRPCServer struct {
    types2.TaskvaultServer
    agent  *Agent
    logger *logrus.Entry
}

func NewGRPCServer(agent *Agent, logger *logrus.Entry) TaskvaultGRPCServer {
    return &GRPCServer{
        agent:  agent,
        logger: logger,
    }
}

func (grpcs *GRPCServer) Serve(lis net.Listener) error {
    grpcServer := grpc.NewServer()
    types2.RegisterTaskvaultServer(grpcServer, grpcs)

    go grpcServer.Serve(lis)

    return nil
}

func Encode(t MessageType, msg interface{}) ([]byte, error) {
    var buf bytes.Buffer
    buf.WriteByte(uint8(t))
    m, err := proto.Marshal(msg.(proto.Message))
    if err != nil {
        return nil, err
    }
    _, err = buf.Write(m)
    return buf.Bytes(), err
}

func (g *GRPCServer) CreateValue(
    ctx context.Context,
    req *types2.CreateValueRequest,
) (*types2.CreateValueResponse, error) {
    defer metrics.MeasureSince([]string{"grpc", "create_value"}, time.Now())
    g.logger.WithFields(
        logrus.Fields{
            "key": req.Key,
            "val": req.Value,
        },
    ).Debug("grpc: Received CreateValue")

    if err := g.agent.applySetPair(
        &types2.Pair{
            Key:   req.Key,
            Value: req.Value,
        },
    ); err != nil {
        return nil, err
    }

    err := g.agent.Store.SetValue(req.Key, req.Value)
    if err != nil {
        return nil, err
    }

    return &types2.CreateValueResponse{
        Key:   req.Key,
        Value: req.Value,
    }, nil
}

func (g *GRPCServer) DeleteValue(
    ctx context.Context,
    req *types2.DeleteValueRequest,
) (*types2.DeleteValueResponse, error) {
    defer metrics.MeasureSince([]string{"grpc", "delete_value"}, time.Now())
    g.logger.WithFields(
        logrus.Fields{
            "key": req.Key,
        },
    ).Debug("grpc: Received DeleteValue")

    cmd, err := Encode(DeletePairType, &req)
    if err != nil {
        return nil, err
    }

    af := g.agent.raft.Apply(cmd, raftTimeout)
    if err := af.Error(); err != nil {
        return nil, err
    }

    err = g.agent.Store.DeleteValue(req.Key)
    if err != nil {
        return nil, err
    }

    res := af.Response()
    resm, ok := res.(*types2.DeleteValueResponse)
    if !ok {
        return nil, fmt.Errorf(
            "grpc: Error wrong response from apply in DeleteValue: %v", res,
        )
    }

    return resm, nil
}

func (g *GRPCServer) GetAllPairs(
    ctx context.Context,
    req *emptypb.Empty,
) (*types2.GetAllPairsResponse, error) {
    defer metrics.MeasureSince([]string{"grpc", "get_all_pairs"}, time.Now())
    g.logger.Debug("grpc: Received GetAllPairs")

    pairs, err := g.agent.Store.GetAllValues()
    if err != nil {
        return nil, err
    }

    p := make([]*types2.Pair, len(pairs))
    for i, pair := range pairs {
        p[i] = &types2.Pair{
            Key:   pair.Key,
            Value: pair.Value,
        }
    }

    return &types2.GetAllPairsResponse{
        Pairs: p,
    }, nil
}

func (g *GRPCServer) GetValue(
    ctx context.Context,
    req *types2.GetValueRequest,
) (*types2.GetValueResponse, error) {
    defer metrics.MeasureSince([]string{"grpc", "get_value"}, time.Now())
    g.logger.WithField("job", req.Key).Debug("grpc: Received GetValue")

    pair, err := g.agent.Store.GetValue(req.Key)
    if err != nil {
        return nil, err
    }

    return &types2.GetValueResponse{
        Value: pair,
    }, nil
}

func (g *GRPCServer) Leave(
    ctx context.Context, req *emptypb.Empty,
) (*emptypb.Empty, error) {
    return req, g.agent.Stop()
}

func (g *GRPCServer) RaftGetConfiguration(
    ctx context.Context,
    req *emptypb.Empty,
) (*types2.RaftGetConfigurationResponse, error) {
    future := g.agent.raft.GetConfiguration()
    if err := future.Error(); err != nil {
        return nil, err
    }

    serverMap := make(map[raft.ServerAddress]serf.Member)
    for _, member := range g.agent.serf.Members() {
        valid, parts := isServer(member)
        if !valid {
            continue
        }

        addr := (&net.TCPAddr{IP: member.Addr, Port: parts.Port}).String()
        serverMap[raft.ServerAddress(addr)] = member
    }

    leader := g.agent.raft.Leader()
    reply := &types2.RaftGetConfigurationResponse{}
    reply.Index = future.Index()
    for _, server := range future.Configuration().Servers {
        node := "(unknown)"
        raftProtocolVersion := "unknown"
        if member, ok := serverMap[server.Address]; ok {
            node = member.Name
            if raftVsn, ok := member.Tags["raft_vsn"]; ok {
                raftProtocolVersion = raftVsn
            }
        }

        entry := &types2.RaftServer{
            Id:           string(server.ID),
            Node:         node,
            Address:      string(server.Address),
            Leader:       server.Address == leader,
            Voter:        server.Suffrage == raft.Voter,
            RaftProtocol: raftProtocolVersion,
        }
        reply.Servers = append(reply.Servers, entry)
    }
    return reply, nil
}

func (g *GRPCServer) RaftRemovePeerByID(
    ctx context.Context,
    req *types2.RaftRemovePeerByIDRequest,
) (*emptypb.Empty, error) {
    panic("unimplemeneted")
}

func (g *GRPCServer) UpdateValue(
    ctx context.Context,
    req *types2.UpdateValueRequest,
) (*types2.UpdateValueResponse, error) {
    panic("unimplemented")
}

