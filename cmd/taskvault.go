package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/danluki/taskvault/logging"
	"github.com/danluki/taskvault/taskvault"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var config = taskvault.DefaultConfig()

var taskvaultCmd = &cobra.Command{
	Use:   "taskvault",
	Short: "Open source distributed job scheduling system",
	Long:  "Task value is a system service that runs scheduled jobs at given intervals or times, just like the unix cron service but distributed in several machines in a cluster. If a machine fails (the leader), a follower will take over and keep running the scheduled jobs without human intervention.",
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
