package taskvault

import (
	"net"
	"time"

	"github.com/hashicorp/raft"
	"go.uber.org/zap"
)

type RaftLayer struct {
	ln     net.Listener
	logger *zap.SugaredLogger
}

var _ raft.StreamLayer = (*RaftLayer)(nil)

func NewRaftLayer(logger *zap.SugaredLogger) *RaftLayer {
	return &RaftLayer{logger: logger}
}

func NewTLSRaftLayer(logger *zap.SugaredLogger) *RaftLayer {
	return &RaftLayer{
		logger: logger,
	}
}

func (t *RaftLayer) Open(l net.Listener) {
	t.ln = l
}

func (t *RaftLayer) Dial(addr raft.ServerAddress, timeout time.Duration) (net.Conn, error) {
	dialer := &net.Dialer{Timeout: timeout}

	var err error
	var conn net.Conn

	conn, err = dialer.Dial("tcp", string(addr))

	return conn, err
}

func (t *RaftLayer) Accept() (net.Conn, error) {
	c, err := t.ln.Accept()
	if err != nil {
		t.logger.Error(err)
	}

	return c, err
}

func (t *RaftLayer) Close() error {
	return t.ln.Close()
}

func (t *RaftLayer) Addr() net.Addr {
	return t.ln.Addr()
}
