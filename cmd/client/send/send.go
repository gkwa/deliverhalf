/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	cmd "github.com/taylormonacelli/deliverhalf/cmd/client"
	common "github.com/taylormonacelli/deliverhalf/cmd/common"
	meta "github.com/taylormonacelli/deliverhalf/cmd/meta"
	sns "github.com/taylormonacelli/deliverhalf/cmd/sns"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("send called")
		logger := common.SetupLogger()
		send(logger)
	},
}

func init() {
	cmd.ClientCmd.AddCommand(sendCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sendCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func send(logger *log.Logger) {
	data := meta.Fetch(logger)
	jsBytes, _ := json.MarshalIndent(data, "", "    ")
	jsonStr := string(jsBytes)
	sns.SendJsonStr(logger, jsonStr)
}
