package taskvault

import (
	"context"
	"time"

	types2 "github.com/danluki/taskvault/pkg/types"
	metrics "github.com/hashicorp/go-metrics"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TaskvaultGRPCClient interface {
	Connect(string) (*grpc.ClientConn, error)
	CreateValue(string, string) (*Pair, error)
	UpdateValue(string, string) (*Pair, error)
	GetValue(string, string) (*Pair, error)
	GetAllValues() ([]Pair, error)
	DeleteValue(string) error
	Leave(string) error
	RaftGetConfiguration(string) (*types2.RaftGetConfigurationResponse, error)
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
func (grpcc *GRPCClient) CreateValue(key string, value string) (*Pair, error) {
	defer metrics.MeasureSince([]string{"grpc", "create_value"}, time.Now())
	var conn *grpc.ClientConn

	addr := grpcc.agent.raft.Leader()

	// Initiate a connection with the server
	conn, err := grpcc.Connect(string(addr))
	if err != nil {
		grpcc.logger.WithError(err).WithFields(
			logrus.Fields{
				"method":      "CreateValue",
				"server_addr": addr,
			},
		).Error("grpc: error dialing.")
		return nil, err
	}
	defer conn.Close()

	d := types2.NewTaskvaultClient(conn)
	resp, err := d.CreateValue(
		context.Background(), &types2.CreateValueRequest{
			Key:   key,
			Value: value,
		},
	)
	if err != nil {
		grpcc.logger.WithError(err).WithFields(
			logrus.Fields{
				"method":      "GetValue",
				"server_addr": addr,
			},
		).Error("grpc: Error calling gRPC method")
		return nil, err
	}

	return &Pair{
		Key:   key,
		Value: resp.Value,
	}, nil
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
func (grpcc *GRPCClient) GetValue(addr, key string) (*Pair, error) {
	defer metrics.MeasureSince([]string{"grpc", "get_value"}, time.Now())
	var conn *grpc.ClientConn

	// Initiate a connection with the server
	conn, err := grpcc.Connect(addr)
	if err != nil {
		grpcc.logger.WithError(err).WithFields(
			logrus.Fields{
				"method":      "GetJob",
				"server_addr": addr,
			},
		).Error("grpc: error dialing.")
		return nil, err
	}
	defer conn.Close()

	d := types2.NewTaskvaultClient(conn)
	resp, err := d.GetValue(
		context.Background(), &types2.GetValueRequest{
			Key: key,
		},
	)
	if err != nil {
		grpcc.logger.WithError(err).WithFields(
			logrus.Fields{
				"method":      "GetValue",
				"server_addr": addr,
			},
		).Error("grpc: Error calling gRPC method")
		return nil, err
	}

	return &Pair{
		Key:   key,
		Value: resp.Value,
	}, nil
}

// Leave implements TaskvaultGRPCClient.
func (grpcc *GRPCClient) Leave(addr string) error {
	var conn *grpc.ClientConn

	// Initiate a connection with the server
	conn, err := grpcc.Connect(addr)
	if err != nil {
		grpcc.logger.WithError(err).WithFields(
			logrus.Fields{
				"method":      "Leave",
				"server_addr": addr,
			},
		).Error("grpc: error dialing.")
		return err
	}
	defer conn.Close()

	// Synchronous call
	d := types2.NewTaskvaultClient(conn)
	_, err = d.Leave(context.Background(), &emptypb.Empty{})
	if err != nil {
		grpcc.logger.WithError(err).WithFields(
			logrus.Fields{
				"method":      "Leave",
				"server_addr": addr,
			},
		).Error("grpc: Error calling gRPC method")
		return err
	}

	return nil
}

// RaftRemovePeerByID implements TaskvaultGRPCClient.
func (grpcc *GRPCClient) RaftRemovePeerByID(addr string, peerID string) error {
	var conn *grpc.ClientConn

	// Initiate a connection with the server
	conn, err := grpcc.Connect(addr)
	if err != nil {
		grpcc.logger.WithError(err).WithFields(
			logrus.Fields{
				"method":      "RaftRemovePeerByID",
				"server_addr": addr,
			},
		).Error("grpc: error dialing.")
		return err
	}
	defer conn.Close()

	// Synchronous call
	d := types2.NewTaskvaultClient(conn)
	_, err = d.RaftRemovePeerByID(
		context.Background(),
		&types2.RaftRemovePeerByIDRequest{Id: peerID},
	)
	if err != nil {
		grpcc.logger.WithError(err).WithFields(
			logrus.Fields{
				"method":      "RaftRemovePeerByID",
				"server_addr": addr,
			},
		).Error("grpc: Error calling gRPC method")
		return err
	}

	return nil
}

// UpdateValue implements TaskvaultGRPCClient.
func (grpcc *GRPCClient) UpdateValue(string, string) (*Pair, error) {
	panic("unimplemented")
}

// RaftGetConfiguration implements TaskvaultGRPCServer.
func (g *GRPCClient) RaftGetConfiguration(
	addr string,
) (*types2.RaftGetConfigurationResponse, error) {
	var conn *grpc.ClientConn

	// Initiate a connection with the server
	conn, err := g.Connect(addr)
	if err != nil {
		g.logger.WithError(err).WithFields(
			logrus.Fields{
				"method":      "RaftGetConfiguration",
				"server_addr": addr,
			},
		).Error("grpc: error dialing.")
		return nil, err
	}
	defer conn.Close()

	// Synchronous call
	d := types2.NewTaskvaultClient(conn)
	res, err := d.RaftGetConfiguration(context.Background(), &emptypb.Empty{})
	if err != nil {
		g.logger.WithError(err).WithFields(
			logrus.Fields{
				"method":      "RaftGetConfiguration",
				"server_addr": addr,
			},
		).Error("grpc: Error calling gRPC method")
		return nil, err
	}

	return res, nil
}
