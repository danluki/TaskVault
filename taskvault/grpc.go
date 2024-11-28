package taskvault

import (
	"bytes"
	"context"
	"net"
	"time"

	"github.com/armon/go-metrics"
	"github.com/danluki/taskvault/types"
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

	// if err := g.agent.apply

	return nil, nil
}

// DeleteJob implements TaskvaultGRPCServer.
// Subtle: this method shadows the method (TaskvaultServer).DeleteJob of GRPCServer.TaskvaultServer.
func (g *GRPCServer) DeleteJob(
	ctx context.Context,
	req *types.DeleteValueRequest,
) (*types.DeleteValueResponse, error) {
	panic("unimplemented")
}

// GetAllPairs implements TaskvaultGRPCServer.
// Subtle: this method shadows the method (TaskvaultServer).GetAllPairs of GRPCServer.TaskvaultServer.
func (g *GRPCServer) GetAllPairs(
	ctx context.Context,
	req *emptypb.Empty,
) (*types.GetAllPairsResponse, error) {
	panic("unimplemented")
}

// GetValue implements TaskvaultGRPCServer.
// Subtle: this method shadows the method (TaskvaultServer).GetValue of GRPCServer.TaskvaultServer.
func (g *GRPCServer) GetValue(
	ctx context.Context,
	req *types.GetValueRequest,
) (*types.GetValueResponse, error) {
	panic("unimplemented")
}

// Leave implements TaskvaultGRPCServer.
// Subtle: this method shadows the method (TaskvaultServer).Leave of GRPCServer.TaskvaultServer.
func (g *GRPCServer) Leave(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	panic("unimplemented")
}

// RaftGetConfiguration implements TaskvaultGRPCServer.
// Subtle: this method shadows the method (TaskvaultServer).RaftGetConfiguration of GRPCServer.TaskvaultServer.
func (g *GRPCServer) RaftGetConfiguration(
	ctx context.Context,
	req *emptypb.Empty,
) (*types.RaftGetConfigurationResponse, error) {
	panic("unimplemented")
}

// RaftRemovePeerByID implements TaskvaultGRPCServer.
// Subtle: this method shadows the method (TaskvaultServer).RaftRemovePeerByID of GRPCServer.TaskvaultServer.
func (g *GRPCServer) RaftRemovePeerByID(
	ctx context.Context,
	req *types.RaftRemovePeerByIDRequest,
) (*emptypb.Empty, error) {
	panic("unimplemented")
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
