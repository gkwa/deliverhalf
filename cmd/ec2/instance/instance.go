/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	cmd "github.com/taylormonacelli/deliverhalf/cmd/ec2"
	log "github.com/taylormonacelli/deliverhalf/cmd/logging"
)

// instanceCmd represents the instance command
var instanceCmd = &cobra.Command{
	Use:   "instance",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger.Traceln("instance called")
	},
}

func init() {
	cmd.Ec2Cmd.AddCommand(instanceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// instanceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// instanceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func describeInstance() {
}
