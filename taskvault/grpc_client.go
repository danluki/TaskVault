package taskvault

import (
	"context"
	"time"

	types2 "github.com/danluki/taskvault/pkg/types"
	metrics "github.com/hashicorp/go-metrics"
	"go.uber.org/zap"
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
	logger  *zap.SugaredLogger
}

func NewGRPCClient(
	dialOpt grpc.DialOption,
	agent *Agent,
	logger *zap.SugaredLogger,
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

func (grpcc *GRPCClient) Connect(addr string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, grpcc.dialOpt...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (grpcc *GRPCClient) CreateValue(key string, value string) (*Pair, error) {
	defer metrics.MeasureSince([]string{"grpc", "create_value"}, time.Now())
	var conn *grpc.ClientConn

	addr := grpcc.agent.raft.Leader()

	conn, err := grpcc.Connect(string(addr))
	if err != nil {
		grpcc.logger.Error("grpc: error dialing",
			zap.Error(err),
			zap.String("method", "CreateValue"),
		)
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
		grpcc.logger.Error("grpc: error calling",
			zap.Error(err),
			zap.String("method", "GetValue"),
		)
		return nil, err
	}

	return &Pair{
		Key:   key,
		Value: resp.Value,
	}, nil
}

func (grpcc *GRPCClient) DeleteValue(string) error {
	panic("unimplemented")
}

func (grpcc *GRPCClient) GetAllValues() ([]Pair, error) {
	panic("unimplemented")
}

func (grpcc *GRPCClient) GetValue(addr, key string) (*Pair, error) {
	defer metrics.MeasureSince([]string{"grpc", "get_value"}, time.Now())
	var conn *grpc.ClientConn

	conn, err := grpcc.Connect(addr)
	if err != nil {
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
		return nil, err
	}

	return &Pair{
		Key:   key,
		Value: resp.Value,
	}, nil
}

func (grpcc *GRPCClient) Leave(addr string) error {
	var conn *grpc.ClientConn

	conn, err := grpcc.Connect(addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	d := types2.NewTaskvaultClient(conn)
	_, err = d.Leave(context.Background(), &emptypb.Empty{})
	if err != nil {
		return err
	}

	return nil
}

func (grpcc *GRPCClient) RaftRemovePeerByID(addr string, peerID string) error {
	var conn *grpc.ClientConn

	conn, err := grpcc.Connect(addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	d := types2.NewTaskvaultClient(conn)
	_, err = d.RaftRemovePeerByID(
		context.Background(),
		&types2.RaftRemovePeerByIDRequest{Id: peerID},
	)
	if err != nil {
		return err
	}

	return nil
}

func (grpcc *GRPCClient) UpdateValue(string, string) (*Pair, error) {
	panic("unimplemented")
}

func (g *GRPCClient) RaftGetConfiguration(
	addr string,
) (*types2.RaftGetConfigurationResponse, error) {
	var conn *grpc.ClientConn

	conn, err := g.Connect(addr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	d := types2.NewTaskvaultClient(conn)
	res, err := d.RaftGetConfiguration(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	return res, nil
}
