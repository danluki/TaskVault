package cmd

import (
	"fmt"

	"github.com/danluki/taskvault/taskvault"
	"github.com/hashicorp/serf/serf"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Long:  `Show the version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Name: %s\n", taskvault.Name)
		fmt.Printf("Version: %s\n", taskvault.Version)
		fmt.Printf("Codename: %s\n", taskvault.Codename)
		fmt.Printf("Agent Protocol: %d (Understands back to: %d)\n",
			serf.ProtocolVersionMax, serf.ProtocolVersionMin)
	},
}

func init() {
	taskvaultCmd.AddCommand(versionCmd)
}
