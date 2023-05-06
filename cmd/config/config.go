/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/taylormonacelli/deliverhalf/cmd"
	common "github.com/taylormonacelli/deliverhalf/cmd/common"
	log "github.com/taylormonacelli/deliverhalf/cmd/logging"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Args:  cobra.OnlyValidArgs,
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger.Trace("config called")
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		os.Exit(1)
		return nil
	},
}

func init() {
	cmd.RootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func reloadConfig() {
	// Read the default configuration file
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file: %s", err)
		return
	}

	// Reload the default configuration file
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file: %s", err)
		return
	}
}

func s3ConfigAbsPath() string {
	bucket := viper.GetString("s3bucket.name")
	s3ConfigPath := viper.GetString("s3bucket.path")

	// Print the object's full path
	path := fmt.Sprintf("s3://%s/%s", bucket, s3ConfigPath)
	return path
}

func showSettings() {
	// Get all configuration settings as a map
	settings := viper.AllSettings()
	common.PrintMap(settings, "")
}
