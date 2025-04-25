package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/spf13/cobra"
)

var keygenCmd = &cobra.Command{
	Use:   "keygen",
	Short: "Generates a new encryption key",
	Long:  `Generates a new encryption key.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		key := make([]byte, 16)
		_, err := rand.Reader.Read(key)
		if err != nil {
			return fmt.Errorf("error reading random data: %s", err)
		}

		fmt.Println(base64.StdEncoding.EncodeToString(key))
		return nil
	},
}

func init() {
	taskvaultCmd.AddCommand(keygenCmd)
}
