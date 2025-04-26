package taskvault

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/hashicorp/raft"
	"go.uber.org/zap"
)

type RaftLayer struct {
	TLSConfig *tls.Config

	ln     net.Listener
	logger *zap.SugaredLogger
}

func NewRaftLayer(logger *zap.SugaredLogger) *RaftLayer {
	return &RaftLayer{logger: logger}
}

func NewTLSRaftLayer(tlsConfig *tls.Config, logger *zap.SugaredLogger) *RaftLayer {
	return &RaftLayer{
		TLSConfig: tlsConfig,
		logger:    logger,
	}
}

func (t *RaftLayer) Open(l net.Listener) error {
	t.ln = l
	return nil
}

func (t *RaftLayer) Dial(addr raft.ServerAddress, timeout time.Duration) (net.Conn, error) {
	dialer := &net.Dialer{Timeout: timeout}

	var err error
	var conn net.Conn
	if t.TLSConfig != nil {
		t.logger.Debug("doing a TLS dial")
		conn, err = tls.DialWithDialer(dialer, "tcp", string(addr), t.TLSConfig)
	} else {
		conn, err = dialer.Dial("tcp", string(addr))
	}

	return conn, err
}

func (t *RaftLayer) Accept() (net.Conn, error) {
	c, err := t.ln.Accept()
	if err != nil {
		fmt.Println("error accepting: ", err.Error())
	}
	return c, err
}

func (t *RaftLayer) Close() error {
	return t.ln.Close()
}

func (t *RaftLayer) Addr() net.Addr {
	return t.ln.Addr()
}
