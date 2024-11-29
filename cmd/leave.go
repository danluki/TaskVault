package cmd

import (
	"github.com/danluki/taskvault/taskvault"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var leaveCmd = &cobra.Command{
	Use:   "leave",
	Short: "Force an agent to leave the cluster",
	Long: `Stop stops an agent, if the agent is a server and is running for election
	stop running for election, if this server was the leader
	this will force the cluster to elect a new leader and start a new scheduler.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		ipa, err := taskvault.ParseSingleIPTemplate(rpcAddr)
		if err != nil {
			return err
		}
		ip = ipa

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var gc taskvault.TaskvaultGRPCClient
		log := logrus.NewEntry(logrus.New())
		gc = taskvault.NewGRPCClient(nil, nil, log)

		if err := gc.Leave(ip); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	taskvaultCmd.AddCommand(leaveCmd)
	leaveCmd.PersistentFlags().
		StringVar(&rpcAddr, "rpc-addr", "{{ GetPrivateIP }}:6868", "gRPC address of the agent")
}
