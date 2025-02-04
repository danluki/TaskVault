package taskvault

import (
	"time"

	"github.com/hashicorp/go-metrics"
	"github.com/hashicorp/go-metrics/prometheus"
)

func initMetrics(a *Agent) error {
	inm := metrics.NewInmemSink(10*time.Second, time.Minute)
	metrics.DefaultInmemSignal(inm)

	var fanout metrics.FanoutSink

	if a.config.EnablePrometheus {
		promSink, err := prometheus.NewPrometheusSink()
		if err != nil {
			return err
		}

		fanout = append(fanout, promSink)
	}

	if len(fanout) > 0 {
		fanout = append(fanout, inm)
		if _, err := metrics.NewGlobal(
			metrics.DefaultConfig("taskvault"), fanout,
		); err != nil {
			return err
		}
	} else {
		if _, err := metrics.NewGlobal(
			metrics.DefaultConfig("taskvault"), inm,
		); err != nil {
			return err
		}
	}

	return nil
}
