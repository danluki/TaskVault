package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/danluki/taskvault/pkg/logging"
	"github.com/danluki/taskvault/taskvault"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var config = taskvault.DefaultConfig()

var rpcAddr string
var ip string

var taskvaultCmd = &cobra.Command{
	Use:   "syncra",
	Short: "Open source distributed core",
	Long:  "Syncra is a open soucre distributed core that will makes your high availability stateful setups possible.",
}

func Execute() {
	if err := taskvaultCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() error {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("taskvault")        // name of config file (without extension)
		viper.AddConfigPath("/etc/taskvault")   // call multiple times to add many search paths
		viper.AddConfigPath("$HOME/.taskvault") // call multiple times to add many search paths
		viper.AddConfigPath("./config")         // call multiple times to add many search paths
	}

	viper.SetEnvPrefix("taskvault")
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv() // read in environment variables that match

	// Add hook to set error logs to stderr and regular logs to stdout
	logrus.AddHook(&logging.LogSplitter{})

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		logrus.WithError(err).Info("No valid config found: Applying default values.")
	}

	if err := viper.Unmarshal(config); err != nil {
		return fmt.Errorf("config: Error unmarshalling config: %s", err)
	}

	cliTags := viper.GetStringSlice("tag")
	var tags map[string]string

	if len(cliTags) > 0 {
		tags, err = UnmarshalTags(cliTags)
		if err != nil {
			return fmt.Errorf("config: Error unmarshalling cli tags: %s", err)
		}
	} else {
		tags = viper.GetStringMapString("tags")
	}

	config.Tags = tags

	taskvault.InitLogger(viper.GetString("log-level"), config.NodeName)

	return nil
}
