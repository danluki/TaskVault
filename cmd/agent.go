package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/danluki/taskvault/taskvault"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ShutdownCh chan (struct{})
var agent *taskvault.Agent

const (
	gracefulTimeout = 3 * time.Hour
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Start a taskvault agent",
	Long:  `Start a taskvault agent. It also runs a web UI if needed.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return initConfig()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return agentRun()
	},
}

func init() {
	taskvaultCmd.AddCommand(agentCmd)

	agentCmd.PersistentFlags().StringVar(
		&cfgFile, "config", "", "config file path",
	)
	agentCmd.Flags().AddFlagSet(taskvault.ConfigFlagSet())
	_ = viper.BindPFlags(agentCmd.Flags())
}

func agentRun() error {
	agent = taskvault.NewAgent(config)
	if err := agent.Start(); err != nil {
		return err
	}

	exit := handleSignals()
	if exit != 0 {
		return fmt.Errorf("exit status: %d", exit)
	}

	return nil
}

func handleSignals() int {
	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	var sig os.Signal
	select {
	case s := <-signalCh:
		sig = s
	case err := <-agent.RetryJoinCh():
		fmt.Println("[ERR] agent: Retry join failed: ", err)
		return 1
	case <-ShutdownCh:
		sig = os.Interrupt
	}
	fmt.Printf("Caught signal: %v", sig)

	if sig != syscall.SIGTERM && sig != os.Interrupt {
		return 1
	}

	log.Info("agent: Gracefully shutting down agent...")
	go func() {
		if err := agent.Stop(); err != nil {
			fmt.Printf("Error: %s", err)
			log.Error(fmt.Sprintf("Error: %s", err))
			return
		}
	}()

	gracefulCh := make(chan struct{})

	time.Sleep(1 * time.Second)

	close(gracefulCh)

	select {
	case <-signalCh:
		return 1
	case <-time.After(gracefulTimeout):
		return 1
	case <-gracefulCh:
		return 0
	}
}

func UnmarshalTags(tags []string) (map[string]string, error) {
	result := make(map[string]string)
	for _, tag := range tags {
		parts := strings.SplitN(tag, "=", 2)
		if len(parts) != 2 || len(parts[0]) == 0 {
			return nil, fmt.Errorf("invalid tag: '%s'", tag)
		}
		result[parts[0]] = parts[1]
	}
	return result, nil
}
