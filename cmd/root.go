package cmd

import (
	"fmt"
	"github.com/gabeduke/kube-ephemeral-pod-reaper/pkg/scout"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "kube-ephemeral-pod-reaper",
	Short: "A tool for managing the lifecycle of ephemeral containers in Kubernetes",
	Long: `kube-ephemeral-pod-reaper is a comprehensive tool designed to enhance Kubernetes' management 
           of ephemeral containers. It provides fine-grained control over the lifecycle of ephemeral 
           containers, ensuring efficient resource utilization and streamlined container management.

           This tool comprises two main components:
           1. Scout: Watches for ephemeral containers and marks them for reaping based on certain criteria.
           2. Reaper: Responsible for the actual deletion of the marked ephemeral containers.

           Both components can be configured and controlled via their respective subcommands, offering 
           a flexible and robust solution for managing ephemeral containers in various Kubernetes environments.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(scout.NewScoutCmd())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".kube-ephemeral-pod-reaper" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".kube-ephemeral-pod-reaper")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
