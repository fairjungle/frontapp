package cmds

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultLogLevel = "info" // logrus: "debug" | "info" | "warning | "error" ...
)

var rootCmd = &cobra.Command{
	Use:   "frontapp",
	Short: "Query tool for frontapp service",
}

func init() {
	rootCmd.PersistentFlags().String("apiToken", "", "API token")
	if err := viper.BindPFlag("apiToken", rootCmd.PersistentFlags().Lookup("apiToken")); err != nil {
		panic(err)
	}

	rootCmd.PersistentFlags().StringP("logLevel", "l", defaultLogLevel, "Log level")
	if err := viper.BindPFlag("logLevel", rootCmd.PersistentFlags().Lookup("logLevel")); err != nil {
		panic(err)
	}

	// setup environment variables
	viper.SetEnvPrefix("frontapp")
	viper.AutomaticEnv()

	// setup config file
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/frontapp/")
	viper.AddConfigPath("$HOME/.frontapp")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(err)
		}
	}
}
