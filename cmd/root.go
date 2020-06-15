package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ory/viper"
)

var rootCmd = &cobra.Command{
	Use:   "ory",
	Short: "The ORY CLI",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.AutomaticEnv()
}
