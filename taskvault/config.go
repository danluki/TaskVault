package taskvault

import (
    "crypto/sha256"
    "encoding/base64"
    "encoding/json"
    "errors"
    "fmt"
    "log"
    "net"
    "os"
    "strconv"
    "strings"
    "time"

    "github.com/hashicorp/go-sockaddr/template"
    flag "github.com/spf13/pflag"
)

type Config struct {
    NodeName string `mapstructure:"node-name"`

    BindAddr string `mapstructure:"bind-addr"`

    HTTPAddr string `mapstructure:"http-addr"`

    // Profile for serf: wan, lan, local
    Profile string

    AdvertiseAddr string `mapstructure:"advertise-addr"`

    Tags map[string]string `mapstructure:"tags"`

    EncryptKey string `mapstructure:"encrypt"`

    StartJoin []string `mapstructure:"join"`

    RetryJoinLAN []string `mapstructure:"retry-join"`

    RetryJoinMaxAttemptsLAN int `mapstructure:"retry-max"`

    RetryJoinIntervalLAN time.Duration `mapstructure:"retry-interval"`

    RPCPort int `mapstructure:"rpc-port"`

    AdvertiseRPCPort int `mapstructure:"advertise-rpc-port"`

    LogLevel string `mapstructure:"log-level"`

    // Doesnt work for now
    Datacenter string

    // Doesnt work for now
    Region string

    Bootstrap bool

    BootstrapExpect int `mapstructure:"bootstrap-expect"`

    DataDir string `mapstructure:"data-dir"`

    DevMode bool

    ReconcileInterval time.Duration

    SerfReconnectTimeout string `mapstructure:"serf-reconnect-timeout"`

    EnablePrometheus bool `mapstructure:"enable-prometheus"`
	
    UI bool
}

const (
    DefaultBindPort      int           = 8946
    DefaultRPCPort       int           = 6868
    DefaultRetryInterval time.Duration = 15 * time.Second
)

var ErrResolvingHost = errors.New("error resolving hostname")

func DefaultConfig() *Config {
    hostname, err := os.Hostname()
    if err != nil {
        log.Panic(err)
    }

    tags := map[string]string{}

    return &Config{
        NodeName: hostname,
        BindAddr: fmt.Sprintf(
            "{{ GetPrivateIP }}:%d", DefaultBindPort,
        ),
        HTTPAddr:             ":8080",
        Profile:              "lan",
        LogLevel:             "info",
        RPCPort:              DefaultRPCPort,
        Tags:                 tags,
        DataDir:              "taskvault.data",
        Datacenter:           "dc1",
        Region:               "global",
        ReconcileInterval:    60 * time.Second,
        SerfReconnectTimeout: "24h",
        UI:                   true,
    }
}

func ConfigFlagSet() *flag.FlagSet {
    c := DefaultConfig()
    cmdFlags := flag.NewFlagSet("agent flagset", flag.ContinueOnError)

    cmdFlags.Bool("server", false, "This node is running in server mode")
    cmdFlags.String(
        "node-name", c.NodeName,
        "Name of this node. Must be unique in the cluster",
    )
    cmdFlags.String(
        "bind-addr", c.BindAddr,
        ``,
    )
    cmdFlags.String(
        "advertise-addr", "",
        ``,
    )
    cmdFlags.String(
        "http-addr", c.HTTPAddr,
        ``,
    )
    cmdFlags.String(
        "profile", c.Profile,
        "",
    )
    cmdFlags.StringSlice(
        "join", []string{},
        "",
    )
    cmdFlags.StringSlice(
        "retry-join", []string{},
        ``,
    )
    cmdFlags.Int(
        "retry-max", 0,
        ``,
    )
    cmdFlags.String(
        "retry-interval", DefaultRetryInterval.String(),
        "",
    )
    cmdFlags.StringSlice(
        "tag", []string{},
        `Tags key=value`,
    )
    cmdFlags.String(
        "encrypt", "",
        "16 bytes value",
    )
    cmdFlags.String(
        "log-level", c.LogLevel,
        "Log level (debug|info|warn|error|fatal|panic)",
    )
    cmdFlags.Int(
        "rpc-port", c.RPCPort,
        ``,
    )
    cmdFlags.Int(
        "advertise-rpc-port", 0,
        "Use the value of rpc-port by default",
    )
    cmdFlags.Int(
        "bootstrap-expect", 0,
        ``,
    )
    cmdFlags.String(
        "data-dir", c.DataDir,
        ``,
    )
    cmdFlags.String(
        "datacenter", c.Datacenter,
        ``,
    )
    cmdFlags.String(
        "region", c.Region,
        ``,
    )
    cmdFlags.String(
        "serf-reconnect-timeout", c.SerfReconnectTimeout,
        ``,
    )
    cmdFlags.Bool(
        "bootstrap", false,
        "Bootstrap the cluster.",
    )
    cmdFlags.Bool(
        "ui", true,
        "",
    )

    cmdFlags.Bool(
        "enable-prometheus", true, "",
    )

    return cmdFlags
}

func (c *Config) normalizeAddrs() error {
    if c.BindAddr != "" {
        ipStr, err := ParseSingleIPTemplate(c.BindAddr)
        if err != nil {
            return fmt.Errorf("bind address resolution failed: %v", err)
        }
        c.BindAddr = ipStr
    }

    if c.HTTPAddr != "" {
        ipStr, err := ParseSingleIPTemplate(c.HTTPAddr)
        if err != nil {
            return fmt.Errorf("HTTP address resolution failed: %v", err)
        }
        c.HTTPAddr = ipStr
    }

    addr, err := normalizeAdvertise(
        c.AdvertiseAddr, c.BindAddr, DefaultBindPort, c.DevMode,
    )
    if err != nil {
        return fmt.Errorf(
            "failed to parse advertise address (%v, %v, %v, %v): %w",
            c.AdvertiseAddr,
            c.BindAddr,
            DefaultBindPort,
            c.DevMode,
            err,
        )
    }
    c.AdvertiseAddr = addr

    return nil
}

func ParseSingleIPTemplate(ipTmpl string) (string, error) {
    out, err := template.Parse(ipTmpl)
    if err != nil {
        return "", fmt.Errorf(
            "unable to parse address template %q: %v", ipTmpl, err,
        )
    }

    ips := strings.Split(out, " ")
    switch len(ips) {
    case 0:
        return "", errors.New("no addresses found, please configure one")
    case 1:
        return ips[0], nil
    default:
        return "", fmt.Errorf(
            "multiple addresses found (%q), please configure one", out,
        )
    }
}

func normalizeAdvertise(
    addr string, bind string, defport int, dev bool,
) (string, error) {
    addr, err := ParseSingleIPTemplate(addr)
    if err != nil {
        return "", fmt.Errorf(
            "Error parsing advertise address template: %v", err,
        )
    }

    if addr != "" {
        _, _, err = net.SplitHostPort(addr)
        if err != nil {
            if !isMissingPort(err) && !isTooManyColons(err) {
                return "", fmt.Errorf(
                    "Error parsing advertise address %q: %v", addr, err,
                )
            }

            return net.JoinHostPort(addr, strconv.Itoa(defport)), nil
        }

        return addr, nil
    }

    ips, err := net.LookupIP(bind)
    if err != nil {
        return "", ErrResolvingHost
    }

    for _, ip := range ips {
        if ip.IsLinkLocalUnicast() || ip.IsGlobalUnicast() {
            return net.JoinHostPort(ip.String(), strconv.Itoa(defport)), nil
        }
        if ip.IsLoopback() {
            if dev {
                return net.JoinHostPort(ip.String(), strconv.Itoa(defport)), nil
            }
            return "", fmt.Errorf(
                "defaulting advertise to localhost is unsafe",
            )
        }
    }

    addr, err = ParseSingleIPTemplate("{{ GetPrivateIP }}")
    if err != nil {
        return "", fmt.Errorf(
            "unable to parse default advertise address: %v", err,
        )
    }
    return net.JoinHostPort(addr, strconv.Itoa(defport)), nil
}

func isMissingPort(err error) bool {
    const missingPort = "missing port in address"
    return err != nil && strings.Contains(err.Error(), missingPort)
}

func isTooManyColons(err error) bool {
    const tooManyColons = "too many colons in address"
    return err != nil && strings.Contains(err.Error(), tooManyColons)
}

func (c *Config) AddrParts(address string) (string, int, error) {
    checkAddr := address

START:
    _, _, err := net.SplitHostPort(checkAddr)
    if ae, ok := err.(*net.AddrError); ok && ae.Err == "missing port in address" {
        checkAddr = fmt.Sprintf("%s:%d", checkAddr, DefaultBindPort)
        goto START
    }
    if err != nil {
        return "", 0, err
    }

    addr, err := net.ResolveTCPAddr("tcp", checkAddr)
    if err != nil {
        return "", 0, err
    }

    return addr.IP.String(), addr.Port, nil
}

func (c *Config) EncryptBytes() ([]byte, error) {
    return base64.StdEncoding.DecodeString(c.EncryptKey)
}

func (c *Config) Hash() (string, error) {
    b, err := json.Marshal(c)
    if err != nil {
        return "", err
    }
    sum := sha256.Sum256(b)
    return base64.StdEncoding.EncodeToString(sum[:]), nil
}
