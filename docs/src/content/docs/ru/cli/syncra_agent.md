---
date: 2025-02-06
title: "syncra agent"
slug: syncra_agent
url: /cli/syncra_agent/
---
## syncra agent

Start a taskvault agent

### Synopsis

Start a taskvault agent
It also runs a web UI.

```
syncra agent [flags]
```

### Options

```
      --advertise-addr string           Address used to advertise to other nodes in the cluster. By default,
                                        the bind address is advertised. The value supports 
                                        go-sockaddr/template format.
      --advertise-rpc-port int          Use the value of rpc-port by default
      --bind-addr string                Specifies which address the agent should bind to for network services, 
                                        including the internal gossip protocol and RPC mechanism. This should be 
                                        specified in IP format, and can be used to easily bind all network services 
                                        to the same address. The value supports go-sockaddr/template format.
                                         (default "{{ GetPrivateIP }}:8946")
      --bootstrap                       Bootstrap the cluster.
      --bootstrap-expect int            Provides the number of expected servers in the datacenter. Either this value 
                                        should not be provided or the value must agree with other servers in the 
                                        cluster. When provided, taskvault waits until the specified number of servers are 
                                        available and then bootstraps the cluster. This allows an initial leader to be 
                                        elected automatically. This flag requires server mode.
      --config string                   config file path
      --data-dir string                 Specifies the directory to use for server-specific data, including the 
                                        replicated log. By default, this is the top-level data-dir, 
                                        like [/var/lib/taskvault] (default "taskvault.data")
      --datacenter string               Specifies the data center of the local agent. All members of a datacenter 
                                        should share a local LAN connection. (default "dc1")
      --enable-prometheus               Enable serving prometheus metrics (default true)
      --encrypt string                  Key for encrypting network traffic. Must be a base64-encoded 16-byte key
  -h, --help                            help for agent
      --http-addr string                Address to bind the UI web server to. Only used when server. The value 
                                        supports go-sockaddr/template format. (default ":8080")
      --join strings                    An initial agent to join with. This flag can be specified multiple times
      --log-level string                Log level (debug|info|warn|error|fatal|panic) (default "info")
      --node-name string                Name of this node. Must be unique in the cluster (default "danlukipc")
      --profile string                  Profile is used to control the timing profiles used (default "lan")
      --raft-multiplier int             An integer multiplier used by servers to scale key Raft timing parameters.
                                        Omitting this value or setting it to 0 uses default timing described below. 
                                        Lower values are used to tighten timing and increase sensitivity while higher 
                                        values relax timings and reduce sensitivity. Tuning this affects the time it 
                                        takes to detect leader failures and to perform leader elections, at the expense 
                                        of requiring more network and CPU resources for better performance. By default, 
                                        taskvault will use a lower-performance timing that's suitable for minimal taskvault 
                                        servers, currently equivalent to setting this to a value of 5 (this default 
                                        may be changed in future versions of taskvault, depending if the target minimum 
                                        server profile changes). Setting this to a value of 1 will configure Raft to 
                                        its highest-performance mode is recommended for production taskvault servers. 
                                        The maximum allowed value is 10. (default 1)
      --region string                   Specifies the region the taskvault agent is a member of. A region typically maps 
                                        to a geographic region, for example us, with potentially multiple zones, which 
                                        map to datacenters such as us-west and us-east (default "global")
      --retry-interval string           Time to wait between join attempts. (default "30s")
      --retry-join strings              Address of an agent to join at start time with retries enabled. 
                                        Can be specified multiple times.
      --retry-max int                   Maximum number of join attempts. Defaults to 0, which will retry indefinitely.
      --rpc-port int                    RPC Port used to communicate with clients. Only used when server. 
                                        The RPC IP Address will be the same as the bind address. (default 6868)
      --serf-reconnect-timeout string   This is the amount of time to attempt to reconnect to a failed node before 
                                        giving up and considering it completely gone. In Kubernetes, you might need 
                                        this to about 5s, because there is no reason to try reconnects for default 
                                        24h value. Also Raft behaves oddly if node is not reaped and returned with 
                                        same ID, but different IP.
                                        Format there: https://golang.org/pkg/time/#ParseDuration (default "24h")
      --server                          This node is running in server mode
      --tag strings                     Tag can be specified multiple times to attach multiple key/value tag pairs 
                                        to the given node, specified as key=value
      --ui                              Enable the web UI on this node. The node must be server. (default true)
```

### SEE ALSO

* [syncra](/cli/syncra/)	 - Open source distributed core

###### Auto generated by spf13/cobra on 6-Feb-2025
