/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/spf13/cobra"
	myec2 "github.com/taylormonacelli/deliverhalf/cmd/ec2"
	lt "github.com/taylormonacelli/deliverhalf/cmd/ec2/launchtemplate"
	log "github.com/taylormonacelli/deliverhalf/cmd/logging"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger.Println("create called")
		testCreateEc2InstanceFromLaunchTemplate()
	},
}

func init() {
	InstanceCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func testCreateEc2InstanceFromLaunchTemplate() {
	region := "us-west-2"
	launchTemplateDataFile := "data/GetLaunchTemplateDataOutput/lt-i-0c47cd895db8040c7.json"
	template, err := lt.CreateLaunchTemplateFromFile(launchTemplateDataFile)
	if err != nil {
		log.Logger.Errorln(err)
	}
	templateID := *template.LaunchTemplateId
	latestVersion := "1"

	createEc2InstanceFromLaunchTemplate(region, templateID, launchTemplateDataFile, latestVersion)
}

func createEc2InstanceFromLaunchTemplate(region string, templateID string, launchTemplateDataFile string, latestVersion string) {
	client, err := myec2.GetEc2Client(region)
	if err != nil {
		log.Logger.Errorln(err)
	}

	minCount := int32(1)
	maxCount := int32(1)

	// Specify the launch template ID and set the desired count to 1
	runInstancesInput := &ec2.RunInstancesInput{
		LaunchTemplate: &types.LaunchTemplateSpecification{
			LaunchTemplateId: aws.String(templateID),
		},
		MinCount: &minCount,
		MaxCount: &maxCount,
	}

	runInstancesOutput, err := client.RunInstances(context.TODO(), runInstancesInput)
	if err != nil {
		log.Logger.Errorln("Failed to run instances:", err)
		return
	}

	for _, instance := range runInstancesOutput.Instances {
		log.Logger.Traceln("Instance ID:", *instance.InstanceId)
	}
}
