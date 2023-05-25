/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	common "github.com/taylormonacelli/deliverhalf/cmd/common"
	lt "github.com/taylormonacelli/deliverhalf/cmd/ec2/launchtemplate"
	log "github.com/taylormonacelli/deliverhalf/cmd/logging"
)

// learn1Cmd represents the test1 command
var learn1Cmd = &cobra.Command{
	Use:   "learn1",
	Short: "Learn AWS API: whats the difference between GetLaunchTemplateDataOutput and CreateLaunchTemplateInput",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger.Trace("test1 called")
		compareGetLaunchTemplateDataOutputToCreateLaunchTemplateInput()
	},
}

func init() {
	lt.LaunchtemplateCmd.AddCommand(learn1Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// test1Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// test1Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func compareGetLaunchTemplateDataOutputToCreateLaunchTemplateInput() {
	f1Name := "data/GetLaunchTemplateDataOutput/lt-i-0c31627ed7b52abcb.json"
	f2Name := "data/CreateLaunchTemplateInput/ct-i-0c31627ed7b52abcb.json"

	f1, err := filepath.Abs(f1Name)
	if err != nil {
		log.Logger.Error(err)
	}

	cltInput, err := lt.CreateLaunchTemplateInput(f1)
	if err != nil {
		log.Logger.Fatalln(err)
	}

	ltName := string("test")
	cltInput.LaunchTemplateName = &ltName

	jsBytes, err := json.MarshalIndent(cltInput, "", "  ")
	if err != nil {
		log.Logger.Errorln(err)
	}
	log.Logger.Debug(string(jsBytes))

	f2, err := filepath.Abs(f2Name)
	if err != nil {
		log.Logger.Error(err)
	}
	common.EnsureParentDirectoryExists(f2)

	file, err := os.Create(f2)
	if err != nil {
		fmt.Println("error creating file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(string(jsBytes))
	if err != nil {
		fmt.Println("error writing to file:", err)
		return
	}
}
