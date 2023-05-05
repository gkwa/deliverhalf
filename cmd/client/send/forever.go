/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"log"
	"time"

	"github.com/spf13/cobra"
	common "github.com/taylormonacelli/deliverhalf/cmd/common"
	meta "github.com/taylormonacelli/deliverhalf/cmd/meta"
	sns "github.com/taylormonacelli/deliverhalf/cmd/sns"
)

var delay time.Duration

// foreverCmd represents the forever command
var foreverCmd = &cobra.Command{
	Use:   "forever",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := common.SetupLogger()
		sendForever(logger)
	},
}

func init() {
	sendCmd.AddCommand(foreverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	duration := 1 * time.Minute
	sendCmd.PersistentFlags().DurationVar(&delay, "delay", duration,
		"Delay command execution for a specified duration")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// foreverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func sendForever(logger *log.Logger) error {
	for {
		data := meta.Fetch(logger)
		jsBytes, _ := json.MarshalIndent(data, "", "    ")
		jsonStr := string(jsBytes)
		sns.SendJsonStr(logger, jsonStr)
		logger.Printf("sleeping %s", delay.String())
		time.Sleep(delay)
	}
}
