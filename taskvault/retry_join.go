package taskvault

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	discover "github.com/hashicorp/go-discover"
	discoverk8s "github.com/hashicorp/go-discover/provider/k8s"
	"go.uber.org/zap"
)

func (a *Agent) retryJoinLAN() {
	r := &retryJoiner{
		cluster:     "LAN",
		addrs:       a.config.RetryJoin,
		maxAttempts: a.config.RetryJoinMaxAttempts,
		interval:    a.config.RetryJoinInterval,
		join:        a.JoinLAN,
	}
	if err := r.retryJoin(a.logger); err != nil {
		a.retryJoinCh <- err
	}
}

type retryJoiner struct {
	cluster string

	addrs []string

	maxAttempts int

	interval time.Duration

	join func([]string) (int, error)
}

func (r *retryJoiner) retryJoin(logger *zap.SugaredLogger) error {
	if len(r.addrs) == 0 {
		return nil
	}

	providers := make(map[string]discover.Provider)
	for k, v := range discover.Providers {
		providers[k] = v
	}
	providers["k8s"] = &discoverk8s.Provider{}

	disco, err := discover.New(
		discover.WithUserAgent(UserAgent()),
		discover.WithProviders(providers),
	)
	if err != nil {
		return err
	}

	logger.Info("agent: Joining cluster...", zap.String("cluster", r.cluster))
	attempt := 0
	for {
		var addrs []string
		var err error

		for _, addr := range r.addrs {
			switch {
			case strings.Contains(addr, "provider="):
				servers, _ := disco.Addrs(
					addr,
					log.New(
						os.Stdout, "",
						log.LstdFlags|log.Lshortfile,
					),
				)

				addrs = append(addrs, servers...)
				logger.Infof(
					"agent: Discovered %s servers: %s", r.cluster,
					strings.Join(servers, " "),
				)

			default:
				ipAddr, err := ParseSingleIPTemplate(addr)
				if err != nil {
					continue
				}
				addrs = append(addrs, ipAddr)
			}
		}

		if len(addrs) > 0 {
			n, err := r.join(addrs)
			if err == nil {
				logger.Infof(
					"agent: Join %s completed. Synced with %d initial agents",
					r.cluster,
					n,
				)
				return nil
			}
		}

		if len(addrs) == 0 {
			err = fmt.Errorf("no servers to join")
		}

		attempt++
		if r.maxAttempts > 0 && attempt > r.maxAttempts {
			return fmt.Errorf(
				"agent: max join %s retry exhausted, exiting", r.cluster,
			)
		}

		logger.Warn("agent: Join failed",
			zap.String("cluster", r.cluster),
			zap.Error(err),
			zap.Duration("retry_interval", r.interval),
		)
		time.Sleep(r.interval)
	}
}
