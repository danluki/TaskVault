package taskvault

import (
	"fmt"
	"net"
	"strconv"

	version "github.com/hashicorp/go-version"
	"github.com/hashicorp/serf/serf"
)

type ServerParts struct {
	Name         string
	ID           string
	Port         int
	Bootstrap    bool
	Expect       int
	RaftVersion  int
	BuildVersion *version.Version
	Addr         net.Addr
	RPCAddr      net.Addr
	Status       serf.MemberStatus
}

func (s *ServerParts) String() string {
	return fmt.Sprintf("%s (Addr: %s)",
		s.Name, s.Addr)
}

func (s *ServerParts) Copy() *ServerParts {
	ns := new(ServerParts)
	*ns = *s
	return ns
}

func UserAgent() string {
	return fmt.Sprintf("taskvault/%s", Version)
}

func toServerPart(m serf.Member) *ServerParts {
	_, bootstrap := m.Tags["bootstrap"]

	expect := 0
	expectStr, ok := m.Tags["expect"]
	var err error
	if ok {
		expect, err = strconv.Atoi(expectStr)
		if err != nil {
			return nil
		}
	}
	if expect == 1 {
		bootstrap = true
	}

	rpcIP := net.ParseIP(m.Tags["rpc_addr"])
	if rpcIP == nil {
		rpcIP = m.Addr
	}

	portStr := m.Tags["port"]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil
	}

	buildVersion, err := version.NewVersion(m.Tags["version"])
	if err != nil {
		buildVersion = &version.Version{}
	}

	parts := &ServerParts{
		Name:         m.Name,
		ID:           m.Name,
		Port:         port,
		Bootstrap:    bootstrap,
		Expect:       expect,
		Addr:         &net.TCPAddr{IP: m.Addr, Port: port},
		RPCAddr:      &net.TCPAddr{IP: rpcIP, Port: port},
		BuildVersion: buildVersion,
		Status:       m.Status,
	}

	return parts
}
