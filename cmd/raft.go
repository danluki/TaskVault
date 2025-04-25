package cmd

import (
	"fmt"

	"github.com/danluki/taskvault/taskvault"
	"github.com/ryanuber/columnize"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var raftCmd = &cobra.Command{
	Use:   "raft [command]",
	Short: "Command to perform some raft operations",
	Long:  ``,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		ipa, err := taskvault.ParseSingleIPTemplate(rpcAddr)
		if err != nil {
			return err
		}
		ip = ipa

		return nil
	},
}

var raftListCmd = &cobra.Command{
	Use:   "list-peers",
	Short: "Command to list raft peers",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log := logrus.NewEntry(logrus.New())
		gc := taskvault.NewGRPCClient(nil, nil, log)

		reply, err := gc.RaftGetConfiguration(ip)
		if err != nil {
			return err
		}

		result := []string{"Node|ID|Address|State|Voter"}
		for _, s := range reply.Servers {
			state := "follower"
			if s.Leader {
				state = "leader"
			}
			result = append(result, fmt.Sprintf("%s|%s|%s|%s|%v",
				s.Node, s.Id, s.Address, state, s.Voter))
		}

		fmt.Println(columnize.SimpleFormat(result))

		return nil
	},
}

var peerID string

var raftRemovePeerCmd = &cobra.Command{
	Use:   "remove-peer",
	Short: "Command to list raft peers",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log := logrus.NewEntry(logrus.New())
		gc := taskvault.NewGRPCClient(nil, nil, log)

		if err := gc.RaftRemovePeerByID(ip, peerID); err != nil {
			return err
		}
		fmt.Println("Peer removed")

		return nil
	},
}

func init() {
	raftCmd.PersistentFlags().
		StringVar(&rpcAddr, "rpc-addr", "{{ GetPrivateIP }}:6868", "gRPC address of the agent.")
	raftRemovePeerCmd.Flags().
		StringVar(&peerID, "peer-id", "", "Remove a taskvault server with the given ID from the Raft configuration.")

	raftCmd.AddCommand(raftListCmd)
	raftCmd.AddCommand(raftRemovePeerCmd)

	taskvaultCmd.AddCommand(raftCmd)
}
