/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/taylormonacelli/deliverhalf/cmd"
	log "github.com/taylormonacelli/deliverhalf/cmd/logging"
)

// WatchdogCmd represents the watchdog command
var WatchdogCmd = &cobra.Command{
	Use:   "watchdog",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger.Traceln("watchdog called")
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
	cmd.RootCmd.AddCommand(WatchdogCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// watchdogCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// watchdogCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
