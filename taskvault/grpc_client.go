package taskvault

import (
	"context"
	"time"

	"github.com/danluki/taskvault/types"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TaskvaultGRPCClient interface {
	Connect(string) (*grpc.ClientConn, error)
	CreateValue(string, string) (*Pair, error)
	UpdateValue(string, string) (*Pair, error)
	GetValue(string) (*Pair, error)
	GetAllValues() ([]Pair, error)
	DeleteValue(string) error
	Leave(string) error
	RaftGetConfiguration(string) (*types.RaftGetConfigurationResponse, error)
	RaftRemovePeerByID(string, string) error
}

type GRPCClient struct {
	dialOpt []grpc.DialOption
	agent   *Agent
	logger  *logrus.Entry
}

func NewGRPCClient(
	dialOpt grpc.DialOption,
	agent *Agent,
	logger *logrus.Entry,
) TaskvaultGRPCClient {
	if dialOpt == nil {
		dialOpt = grpc.WithInsecure()
	}
	return &GRPCClient{
		dialOpt: []grpc.DialOption{
			dialOpt,
			grpc.WithBlock(),
		},
		agent:  agent,
		logger: logger,
	}
}

// Connect dialing to a gRPC server
func (grpcc *GRPCClient) Connect(addr string) (*grpc.ClientConn, error) {
	// Initiate a connection with the server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, grpcc.dialOpt...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// CreateValue implements TaskvaultGRPCClient.
func (grpcc *GRPCClient) CreateValue(string, string) (*Pair, error) {
	panic("unimplemented")
}

// DeleteValue implements TaskvaultGRPCClient.
func (grpcc *GRPCClient) DeleteValue(string) error {
	panic("unimplemented")
}

// GetAllValues implements TaskvaultGRPCClient.
func (grpcc *GRPCClient) GetAllValues() ([]Pair, error) {
	panic("unimplemented")
}

// GetValue implements TaskvaultGRPCClient.
func (grpcc *GRPCClient) GetValue(string) (*Pair, error) {
	panic("unimplemented")
}

// Leave implements TaskvaultGRPCClient.
func (grpcc *GRPCClient) Leave(string) error {
	panic("unimplemented")
}

// RaftRemovePeerByID implements TaskvaultGRPCClient.
func (grpcc *GRPCClient) RaftRemovePeerByID(string, string) error {
	panic("unimplemented")
}

// UpdateValue implements TaskvaultGRPCClient.
func (grpcc *GRPCClient) UpdateValue(string, string) (*Pair, error) {
	panic("unimplemented")
}


// RaftGetConfiguration implements TaskvaultGRPCServer.
func (g *GRPCClient) RaftGetConfiguration(
	addr string,
) (*types.RaftGetConfigurationResponse, error) {
	var conn *grpc.ClientConn

	// Initiate a connection with the server
	conn, err := g.Connect(addr)
	if err != nil {
		g.logger.WithError(err).WithFields(logrus.Fields{
			"method":      "RaftGetConfiguration",
			"server_addr": addr,
		}).Error("grpc: error dialing.")
		return nil, err
	}
	defer conn.Close()

	// Synchronous call
	d := types.NewTaskvaultClient(conn)
	res, err := d.RaftGetConfiguration(context.Background(), &emptypb.Empty{})
	if err != nil {
		g.logger.WithError(err).WithFields(logrus.Fields{
			"method":      "RaftGetConfiguration",
			"server_addr": addr,
		}).Error("grpc: Error calling gRPC method")
		return nil, err
	}

	return res, nil
}
