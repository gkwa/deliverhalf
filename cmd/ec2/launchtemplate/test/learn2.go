/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	lt "github.com/taylormonacelli/deliverhalf/cmd/ec2/launchtemplate"
	log "github.com/taylormonacelli/deliverhalf/cmd/logging"
)

// learn2Cmd represents the learn2 command
var learn2Cmd = &cobra.Command{
	Use:   "learn2",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger.Trace("learn2 called")
	},
}

func init() {
	lt.LaunchtemplateCmd.AddCommand(learn2Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// learn2Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// learn2Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
