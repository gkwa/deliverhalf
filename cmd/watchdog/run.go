/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	db "github.com/taylormonacelli/deliverhalf/cmd/db"
	sns "github.com/taylormonacelli/deliverhalf/cmd/sns"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	WatchdogCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func listenForSnsMessages(c chan<- string) {
	// simulate a long running task
	// time.Sleep(3 * time.Second)
	region := viper.GetString("sns.region")
	sns.GetIdentityDocFromSNS(region)
	c <- "longRunningFunc completed"
}

func run() {
	c := make(chan string)

	go listenForSnsMessages(c)
	db.Test2()

	// do other work here while listenForSnsMessages runs in the background
	fmt.Println("Doing other work...")

	// wait for the listenForSnsMessages to complete and print the result
	result := <-c
	fmt.Println(result)
}
