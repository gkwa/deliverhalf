/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/spf13/cobra"
	myec2 "github.com/taylormonacelli/deliverhalf/cmd/ec2"
	volume "github.com/taylormonacelli/deliverhalf/cmd/ec2/volume"
	log "github.com/taylormonacelli/deliverhalf/cmd/logging"
)

// test2Cmd represents the test2 command
var test2Cmd = &cobra.Command{
	Use:   "test2",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger.Traceln("test2 called")
		test2()
	},
}

func init() {
	volume.VolumeCmd.AddCommand(test2Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// test2Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// test2Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func test2() {
	region := "us-west-2"
	svc, err := myec2.GetEc2Client(region)
	if err != nil {
		log.Logger.Fatal(err)
	}

	instanceID := "i-0e602b7a3b4c5299b"

	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	}

	resp, err := svc.DescribeInstances(context.TODO(), input)
	if err != nil {
		log.Logger.Fatal(err)
	}

	// Marshal the response to indented JSON
	jsonBytes, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonBytes))
}
