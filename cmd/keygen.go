package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var keygenCmd = &cobra.Command{
	Use:   "keygen",
	Short: "Generates a new encryption key",
	Long:  `Generates a new encryption key.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		key := make([]byte, 16)
		n, err := rand.Reader.Read(key)
		if err != nil {
			return fmt.Errorf("error reading random data: %s", err)
		}
		if n != 16 {
			return errors.New("couldn't read enough entropy. Generate more entropy")
		}

		fmt.Println(base64.StdEncoding.EncodeToString(key))
		return nil
	},
}

func init() {
	taskvaultCmd.AddCommand(keygenCmd)
}
