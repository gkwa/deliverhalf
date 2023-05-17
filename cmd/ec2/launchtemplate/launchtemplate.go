/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"text/template"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	common "github.com/taylormonacelli/deliverhalf/cmd/common"
	cmd "github.com/taylormonacelli/deliverhalf/cmd/ec2"
	log "github.com/taylormonacelli/deliverhalf/cmd/logging"
)

// LaunchtemplateCmd represents the launchtemplate command
var LaunchtemplateCmd = &cobra.Command{
	Use:   "launchtemplate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
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
	cmd.Ec2Cmd.AddCommand(LaunchtemplateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// launchtemplateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// launchtemplateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func createUserData() (string, error) {
	sshPublicKeys := viper.GetStringSlice("ssh_public_keys")

	tmpl := template.Must(template.New("sshKeys").Parse(`#!/usr/bin/env bash
mkdir -p /root/.ssh
{{range .}}
echo {{.}} >>/root/.ssh/authorized_keys
{{- end}}`))

	var tplOutput bytes.Buffer
	if err := tmpl.Execute(&tplOutput, sshPublicKeys); err != nil {
		return "", nil
	}

	log.Logger.Tracef("userdata script: %s", tplOutput.String())
	return tplOutput.String(), nil
}

func genLaunchTemplateFromInstanceId(region string, instanceID string, ltFname string) {
	// Create a new AWS SDK config with default options
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Logger.Traceln("failed to load SDK config:", err)
		os.Exit(1)
	}

	// Create a new EC2 client
	client := ec2.NewFromConfig(cfg)

	// Retrieve the LaunchTemplateData and write it to a file
	resp, err := getLaunchTemplateDataFromInstanceId(context.Background(), client, instanceID)
	if err != nil {
		log.Logger.Traceln("failed to get LaunchTemplateData:", err)
		os.Exit(1)
	}
	err = writeLaunchTemplateDataToFile(resp, ltFname)
	if err != nil {
		log.Logger.Traceln("failed to write LaunchTemplateData to file:", err)
		os.Exit(1)
	}
}

func getInstanceMap(client *ec2.Client) (map[string]string, error) {
	// Query EC2 instances
	input := &ec2.DescribeInstancesInput{}
	result, err := client.DescribeInstances(context.Background(), input)
	if err != nil {
		return nil, err
	}

	// Create map of instance IDs to instance names
	instances := make(map[string]string)
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			if instance.State.Name != "terminated" && instance.State.Name != "shutting-down" {
				instanceName := ""
				for _, tag := range instance.Tags {
					if *tag.Key == "Name" {
						instanceName = *tag.Value
						break
					}
				}
				instances[*instance.InstanceId] = instanceName
			}
		}
	}
	return instances, nil
}

func getInstanceList(region string, client *ec2.Client) ([]EC2Instance, error) {
	// Query EC2 instances
	input := &ec2.DescribeInstancesInput{}
	result, err := client.DescribeInstances(context.Background(), input)
	if err != nil {
		return nil, err
	}

	// Create slice of instances
	instances := []EC2Instance{}
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			if instance.State.Name != "terminated" && instance.State.Name != "shutting-down" {
				instanceName := ""
				for _, tag := range instance.Tags {
					if *tag.Key == "Name" {
						instanceName = *tag.Value
						break
					}
				}
				instances = append(instances, EC2Instance{InstanceId: *instance.InstanceId, InstanceName: instanceName})
			}
		}
	}
	return instances, nil
}

func genLaunchTemplateFileAbsPath(instancId string) string {
	dir, err := os.Getwd()
	if err != nil {
		log.Logger.Fatalln(err)
	}
	subdir := "data"
	fname := "lt-" + instancId + ".json"

	fullPath := filepath.Join(dir, subdir, fname)
	return fullPath
}

func getBasedirectoryFromPath(filePath string) string {
	baseDir := filepath.Base(filepath.Dir(filePath))
	return baseDir
}

func getAllAwsRegions() []types.Region {
	// Load the AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-west-2"))
	if err != nil {
		log.Logger.Traceln("failed to load AWS SDK config:", err)
		return []types.Region{}
	}

	// Create an EC2 client using the loaded config
	client := ec2.NewFromConfig(cfg)

	// Get a list of all AWS regions
	resp, err := client.DescribeRegions(context.Background(), nil)
	if err != nil {
		log.Logger.Traceln("failed to describe AWS regions:", err)
		return []types.Region{}
	}

	// Create an empty slice of types.Region
	regions := make([]types.Region, 0, len(resp.Regions))
	regions = append(regions, resp.Regions...)

	// Print the region names
	for _, region := range regions {
		log.Logger.Traceln(*region.RegionName)
	}
	return regions
}

func genLaunchTemplatesForAllEc2InstancesInregion(region string) {
	// Create AWS SDK config with default options
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Logger.Fatalln(err)
	}

	// Create EC2 client
	client := ec2.NewFromConfig(cfg)

	// Get instance ID to name map
	instanceMap, err := getInstanceMap(client)
	if err != nil {
		log.Logger.Fatalln(err)
	}

	// Print instance IDs and names
	for id, name := range instanceMap {
		log.Logger.Tracef("Instance ID: %s, Instance Name: %s", id, name)
	}

	// fetch templates locally if not i don't have it
	for id, name := range instanceMap {
		ltPath := genLaunchTemplateFileAbsPath(id)
		dir := getBasedirectoryFromPath(ltPath)
		common.CreateDirectory(dir)
		if common.FileExists(ltPath) {
			log.Logger.Tracef("skipping %s because %s exists",
				name, ltPath)
			continue
		}
		genLaunchTemplateFromInstanceId(region, id, ltPath)
	}
}
