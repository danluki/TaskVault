package taskvault

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"testing"
	"time"

	"github.com/hashicorp/serf/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	logLevel = "info"
)

func TestEncrypt(t *testing.T) {
	dir, err := ioutil.TempDir("", "taskvault-test")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	ip1, returnFn1 := testutil.TakeIP()
	defer returnFn1()

	c := DefaultConfig()
	c.BindAddr = ip1.String()
	c.NodeName = "test1"
	c.EncryptKey = "kPpdjphiipNSsjd4QHWbkA=="
	c.LogLevel = logLevel
	c.DevMode = true
	c.DataDir = dir

	a := NewAgent(c)
	_ = a.Start()

	time.Sleep(2 * time.Second)

	assert.True(t, a.serf.EncryptionEnabled())
	_ = a.Stop()
}

func Test_advertiseRPCAddr(t *testing.T) {
	dir, err := ioutil.TempDir("", "taskvault-test")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	ip1, returnFn1 := testutil.TakeIP()
	defer returnFn1()
	a1Addr := ip1.String()

	c := DefaultConfig()
	c.BindAddr = a1Addr + ":5000"
	c.AdvertiseAddr = "8.8.8.8"
	c.NodeName = "test1"
	c.LogLevel = logLevel
	c.DevMode = true
	c.DataDir = dir

	a := NewAgent(c)
	_ = a.Start()

	time.Sleep(2 * time.Second)

	advertiseRPCAddr := a.advertiseRPCAddr()
	exRPCAddr := "8.8.8.8:6868"

	assert.Equal(t, exRPCAddr, advertiseRPCAddr)

	_ = a.Stop()
}

func Test_bindRPCAddr(t *testing.T) {
	dir, err := ioutil.TempDir("", "taskvault-test")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	ip1, returnFn1 := testutil.TakeIP()
	defer returnFn1()
	a1Addr := ip1.String()

	c := DefaultConfig()
	c.BindAddr = a1Addr + ":5000"
	c.NodeName = "test1"
	c.LogLevel = logLevel
	c.DevMode = true
	c.DataDir = dir
	c.AdvertiseAddr = c.BindAddr

	a := NewAgent(c)
	err = a.Start()
	require.NoError(t, err)

	time.Sleep(2 * time.Second)

	bindRPCAddr := a.bindRPCAddr()
	exRPCAddr := a1Addr + ":6868"
	assert.Equal(t, exRPCAddr, bindRPCAddr)
}

func TestAgentConfig(t *testing.T) {
	dir, err := ioutil.TempDir("", "taskvault-test")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	ip1, returnFn1 := testutil.TakeIP()
	defer returnFn1()
	advAddr := ip1.String()

	ip2, returnFn2 := testutil.TakeIP()
	defer returnFn2()

	c := DefaultConfig()
	c.BindAddr = ip2.String()
	c.AdvertiseAddr = advAddr
	c.LogLevel = logLevel
	c.DataDir = dir
	c.DevMode = true

	a := NewAgent(c)
	_ = a.Start()

	time.Sleep(2 * time.Second)

	assert.NotEqual(t, a.config.AdvertiseAddr, a.config.BindAddr)
	assert.NotEmpty(t, a.config.AdvertiseAddr)
	assert.Equal(t, advAddr+":8946", a.config.AdvertiseAddr)

	_ = a.Stop()
}

func TestAgentCommand_testElection(t *testing.T) {
	dir, err := ioutil.TempDir("", "taskvault-test")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	a1Name := "test1"
	a2Name := "test2"
	ip1, returnFn1 := testutil.TakeIP()
	a1Addr := ip1.String()
	defer returnFn1()
	ip2, returnFn2 := testutil.TakeIP()
	a2Addr := ip2.String()
	defer returnFn2()

	c := DefaultConfig()
	c.BindAddr = a1Addr
	c.StartJoin = []string{a2Addr}
	c.NodeName = a1Name
	c.LogLevel = logLevel
	c.BootstrapExpect = 3
	c.DevMode = true
	c.HTTPAddr = ":8080"
	c.DataDir = dir

	a1 := NewAgent(c)
	err = a1.Start()
	require.NoError(t, err)

	time.Sleep(2 * time.Second)

	if a1.IsLeader() {
		m, err := a1.leaderMember()
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%s is the current leader", m.Name)
		assert.Equal(t, a1Name, m.Name)
	}

	c = DefaultConfig()
	c.BindAddr = a2Addr
	c.StartJoin = []string{a1Addr + ":8946"}
	c.NodeName = a2Name
	c.LogLevel = logLevel
	c.BootstrapExpect = 3
	c.DevMode = true
	c.HTTPAddr = ":8081"
	c.DataDir = dir

	a2 := NewAgent(c)
	err = a2.Start()
	require.NoError(t, err)

	time.Sleep(2 * time.Second)

	c = DefaultConfig()
	ip3, returnFn3 := testutil.TakeIP()
	defer returnFn3()
	c.BindAddr = ip3.String()
	c.StartJoin = []string{a1Addr + ":8946"}
	c.NodeName = "test3"
	c.LogLevel = logLevel
	c.BootstrapExpect = 3
	c.HTTPAddr = ":8082"
	c.DevMode = true
	c.DataDir = dir

	a3 := NewAgent(c)
	err = a3.Start()
	require.NoError(t, err)

	time.Sleep(2 * time.Second)

	_ = a1.Stop()

	time.Sleep(2 * time.Second)
	assert.True(t, (a2.IsLeader() || a3.IsLeader()))
	log.Println(a3.IsLeader())

	_ = a2.Stop()
	_ = a3.Stop()
}

func TestAgents(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	dir, err := ioutil.TempDir("", "taskvault-test")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	a1Name := "test1"
	a2Name := "test2"
	ip1, returnFn1 := testutil.TakeIP()
	a1Addr := ip1.String()
	defer returnFn1()
	ip2, returnFn2 := testutil.TakeIP()
	a2Addr := ip2.String()
	defer returnFn2()

	c := DefaultConfig()
	c.BindAddr = a1Addr
	c.StartJoin = []string{a2Addr}
	c.NodeName = a1Name
	c.LogLevel = logLevel
	c.BootstrapExpect = 3
	c.DevMode = true
	c.HTTPAddr = ":8080"
	c.DataDir = dir

	a1 := NewAgent(c)
	err = a1.Start()
	require.NoError(t, err)

	time.Sleep(2 * time.Second)

	if a1.IsLeader() {
		m, err := a1.leaderMember()
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%s is the current leader", m.Name)
		assert.Equal(t, a1Name, m.Name)
	}

	c = DefaultConfig()
	c.BindAddr = a2Addr
	c.StartJoin = []string{a1Addr + ":8946"}
	c.NodeName = a2Name
	c.LogLevel = logLevel
	c.BootstrapExpect = 3
	c.DevMode = true
	c.HTTPAddr = ":8081"
	c.DataDir = dir

	a2 := NewAgent(c)
	err = a2.Start()
	require.NoError(t, err)

	time.Sleep(2 * time.Second)

	c = DefaultConfig()
	ip3, returnFn3 := testutil.TakeIP()
	defer returnFn3()
	c.BindAddr = ip3.String()
	c.StartJoin = []string{a1Addr + ":8946"}
	c.NodeName = "test3"
	c.LogLevel = logLevel
	c.BootstrapExpect = 3
	c.HTTPAddr = ":8082"
	c.DevMode = true
	c.DataDir = dir

	a3 := NewAgent(c)
	err = a3.Start()
	require.NoError(t, err)

	time.Sleep(2 * time.Second)

	_ = a1.Stop()

	time.Sleep(2 * time.Second)
	assert.True(t, (a2.IsLeader() || a3.IsLeader()))
	log.Println(a3.IsLeader())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
}
