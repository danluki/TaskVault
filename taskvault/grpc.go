package taskvault

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"time"

	"github.com/armon/go-metrics"
	"github.com/danluki/taskvault/types"
	"github.com/hashicorp/raft"
	"github.com/hashicorp/serf/serf"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

// TaskvaultGRPCServer defines the basics that a gRPC server should implement.
type TaskvaultGRPCServer interface {
	types.TaskvaultServer
	Serve(net.Listener) error
}

type GRPCServer struct {
	types.TaskvaultServer
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
	types.RegisterTaskvaultServer(grpcServer, grpcs)

	go grpcServer.Serve(lis)

	return nil
}

// Encode is used to encode a Protoc object with type prefix
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

// CreateValue implements TaskvaultGRPCServer.
// Subtle: this method shadows the method (TaskvaultServer).CreateValue of GRPCServer.TaskvaultServer.
func (g *GRPCServer) CreateValue(
	ctx context.Context,
	req *types.CreateValueRequest,
) (*types.CreateValueResponse, error) {
	defer metrics.MeasureSince([]string{"grpc", "create_value"}, time.Now())
	g.logger.WithFields(logrus.Fields{
		"key": req.Key,
		"val": req.Value,
	}).Debug("grpc: Received CreateValue")

	if err := g.agent.applySetPair(&types.Pair{
		Key:   req.Key,
		Value: req.Value,
	}); err != nil {
		return nil, err
	}

	err := g.agent.Store.SetValue(req.Key, req.Value)
	if err != nil {
		return nil, err
	}

	return &types.CreateValueResponse{
		Key:   req.Key,
		Value: req.Value,
	}, nil
}

// DeleteJob implements TaskvaultGRPCServer.
// Subtle: this method shadows the method (TaskvaultServer).DeleteJob of GRPCServer.TaskvaultServer.
func (g *GRPCServer) DeleteValue(
	ctx context.Context,
	req *types.DeleteValueRequest,
) (*types.DeleteValueResponse, error) {
	defer metrics.MeasureSince([]string{"grpc", "delete_value"}, time.Now())
	g.logger.WithFields(logrus.Fields{
		"key": req.Key,
	}).Debug("grpc: Received DeleteValue")

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
	resm, ok := res.(*types.DeleteValueResponse)
	if !ok {
		return nil, fmt.Errorf("grpc: Error wrong response from apply in DeleteValue: %v", res)
	}

	return resm, nil
}

// GetAllPairs implements TaskvaultGRPCServer.
// Subtle: this method shadows the method (TaskvaultServer).GetAllPairs of GRPCServer.TaskvaultServer.
func (g *GRPCServer) GetAllPairs(
	ctx context.Context,
	req *emptypb.Empty,
) (*types.GetAllPairsResponse, error) {
	defer metrics.MeasureSince([]string{"grpc", "get_all_pairs"}, time.Now())
	g.logger.Debug("grpc: Received GetAllPairs")

	pairs, err := g.agent.Store.GetAllValues()
	if err != nil {
		return nil, err
	}

	p := make([]*types.Pair, len(pairs))
	for i, pair := range pairs {
		p[i] = &types.Pair{
			Key:   pair.Key,
			Value: pair.Value,
		}
	}

	return &types.GetAllPairsResponse{
		Pairs: p,
	}, nil
}

// GetValue implements TaskvaultGRPCServer.
// Subtle: this method shadows the method (TaskvaultServer).GetValue of GRPCServer.TaskvaultServer.
func (g *GRPCServer) GetValue(
	ctx context.Context,
	req *types.GetValueRequest,
) (*types.GetValueResponse, error) {
	defer metrics.MeasureSince([]string{"grpc", "get_value"}, time.Now())
	g.logger.WithField("job", req.Key).Debug("grpc: Received GetValue")

	fmt.Println(req.Key)
	pair, err := g.agent.Store.GetValue(req.Key)
	if err != nil {
		return nil, err
	}

	return &types.GetValueResponse{
		Value: pair,
	}, nil
}

// Leave implements TaskvaultGRPCServer.
// Subtle: this method shadows the method (TaskvaultServer).Leave of GRPCServer.TaskvaultServer.
func (g *GRPCServer) Leave(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return req, g.agent.Stop()
}

// RaftGetConfiguration implements TaskvaultGRPCServer.
// Subtle: this method shadows the method (TaskvaultServer).RaftGetConfiguration of GRPCServer.TaskvaultServer.
func (g *GRPCServer) RaftGetConfiguration(
	ctx context.Context,
	req *emptypb.Empty,
) (*types.RaftGetConfigurationResponse, error) {
	// We can't fetch the leader and the configuration atomically with
	// the current Raft API.
	future := g.agent.raft.GetConfiguration()
	if err := future.Error(); err != nil {
		return nil, err
	}

	// Index the information about the servers.
	serverMap := make(map[raft.ServerAddress]serf.Member)
	for _, member := range g.agent.serf.Members() {
		valid, parts := isServer(member)
		if !valid {
			continue
		}

		addr := (&net.TCPAddr{IP: member.Addr, Port: parts.Port}).String()
		serverMap[raft.ServerAddress(addr)] = member
	}

	// Fill out the reply.
	leader := g.agent.raft.Leader()
	reply := &types.RaftGetConfigurationResponse{}
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

		entry := &types.RaftServer{
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

// RaftRemovePeerByID implements TaskvaultGRPCServer.
// Subtle: this method shadows the method (TaskvaultServer).RaftRemovePeerByID of GRPCServer.TaskvaultServer.
func (g *GRPCServer) RaftRemovePeerByID(
	ctx context.Context,
	req *types.RaftRemovePeerByIDRequest,
) (*emptypb.Empty, error) {
	// Since this is an operation designed for humans to use, we will return
	// an error if the supplied id isn't among the peers since it's
	// likely they screwed up.
	{
		future := g.agent.raft.GetConfiguration()
		if err := future.Error(); err != nil {
			return nil, err
		}
		for _, s := range future.Configuration().Servers {
			if s.ID == raft.ServerID(req.Id) {
				goto REMOVE
			}
		}
		return nil, fmt.Errorf("id %q was not found in the Raft configuration", req.Id)
	}

REMOVE:
	// The Raft library itself will prevent various forms of foot-shooting,
	// like making a configuration with no voters. Some consideration was
	// given here to adding more checks, but it was decided to make this as
	// low-level and direct as possible. We've got ACL coverage to lock this
	// down, and if you are an operator, it's assumed you know what you are
	// doing if you are calling this. If you remove a peer that's known to
	// Serf, for example, it will come back when the leader does a reconcile
	// pass.
	future := g.agent.raft.RemoveServer(raft.ServerID(req.Id), 0, 0)
	if err := future.Error(); err != nil {
		g.logger.WithError(err).WithField("peer", req.Id).Warn("failed to remove Raft peer")
		return nil, err
	}

	g.logger.WithField("peer", req.Id).Warn("removed Raft peer")
	return new(emptypb.Empty), nil
}

// UpdateValue implements TaskvaultGRPCServer.
// Subtle: this method shadows the method (TaskvaultServer).UpdateValue of GRPCServer.TaskvaultServer.
func (g *GRPCServer) UpdateValue(
	ctx context.Context,
	req *types.UpdateValueRequest,
) (*types.UpdateValueResponse, error) {
	panic("unimplemented")
}

// mustEmbedUnimplementedTaskvaultServer implements TaskvaultGRPCServer.
// Subtle: this method shadows the method (TaskvaultServer).mustEmbedUnimplementedTaskvaultServer of GRPCServer.TaskvaultServer.
func (g *GRPCServer) mustEmbedUnimplementedTaskvaultServer() {
	panic("unimplemented")
}
