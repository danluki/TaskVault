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
		viper.SetConfigName("taskvault") 
		viper.AddConfigPath("/etc/taskvault") 
		viper.AddConfigPath("$HOME/.taskvault") 
		viper.AddConfigPath("./config") 
	}

	viper.SetEnvPrefix("taskvault")
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	logrus.AddHook(&logging.LogSplitter{})

	err := viper.ReadInConfig()
	if err != nil {             
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
