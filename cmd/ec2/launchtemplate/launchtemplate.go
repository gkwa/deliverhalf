/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	humanize "github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	common "github.com/taylormonacelli/deliverhalf/cmd/common"
	mydb "github.com/taylormonacelli/deliverhalf/cmd/db"
	cmd "github.com/taylormonacelli/deliverhalf/cmd/ec2"
	myec2 "github.com/taylormonacelli/deliverhalf/cmd/ec2"
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

func GenLaunchTemplateFromInstanceId(region string, instanceID string, ltFname string) (*ec2.GetLaunchTemplateDataOutput, error) {
	db, err := mydb.OpenDB("test.db")
	if err != nil {
		log.Logger.Fatalf("failed to connect to database: %v", err)
	}

	var count int64
	if err := db.Model(&ExtendedGetLaunchTemplateDataOutput{}).Count(&count).Error; err != nil {
		log.Logger.Fatalln(err)
	}

	// don't fetch more than 1 per hour
	since := time.Now().Add(-1 * time.Hour)

	// get most recent template
	var templates []ExtendedGetLaunchTemplateDataOutput
	query := "created_at >= ? and instance_id = ?"
	if err := db.Where(query, since, instanceID).Find(&templates).Error; err != nil {
		log.Logger.Trace("query matchd no results")
		return &ec2.GetLaunchTemplateDataOutput{}, err
	}

	var items []string
	for _, tpl := range templates {
		item := fmt.Sprintf("%s created %s (%s)", tpl.InstanceId, tpl.CreatedAt, humanize.Time(tpl.CreatedAt))
		items = append(items, item)
	}
	log.Logger.Tracef("found %d of %d total templates matching for instance %s",
		count, len(templates), strings.Join(items, ", "))

	// for debug I do want the json file

	// // reduce aws api usage
	// if len(templates) > 0 {
	// 	tpl := templates[len(templates)-1]
	// 	var x1 *ec2.GetLaunchTemplateDataOutput
	// 	js := tpl.LaunchTemplateDataJsonStr
	// 	json.Unmarshal([]byte(js), &x1)
	// 	return x1, nil
	// }

	client, err := myec2.GetEc2Client(region)
	if err != nil {
		log.Logger.Errorln(err)
	}

	resp, err := getLaunchTemplateDataFromInstanceId(context.Background(), client, instanceID)
	if err != nil {
		log.Logger.Errorf("failed to get LaunchTemplateData: %v", err)
	}

	err = writeLaunchTemplateDataToFile(resp, ltFname)
	if err != nil {
		log.Logger.Errorf("failed to write LaunchTemplateData to file %s: %v", ltFname, err)
		return nil, err
	}

	err = writeLaunchTemplateDataForInstanceIdToDB(resp, instanceID)
	if err != nil {
		log.Logger.Errorf("failed to write LaunchTemplateData to db: %v", err)
	}

	return resp, nil
}

func writeLaunchTemplateDataForInstanceIdToDB(resp *ec2.GetLaunchTemplateDataOutput, instancdId string) error {
	jsonData, err := json.Marshal(resp)
	if err != nil {
		log.Logger.Errorf("failed to serialize launchtemplatedataoutput for instance %s", instancdId)
	}

	db, err := mydb.OpenDB("test.db")
	if err != nil {
		log.Logger.Errorf("failed to connect to database: %v", err)
	}

	db.Create(&ExtendedGetLaunchTemplateDataOutput{
		InstanceId:                instancdId,
		LaunchTemplateDataJsonStr: string(jsonData),
	})
	return err
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
	subdir2 := "GetLaunchTemplateDataOutput"
	fname := "lt-" + instancId + ".json"

	fullPath := filepath.Join(dir, subdir, subdir2, fname)
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
	client, err := myec2.GetEc2Client(region)
	if err != nil {
		log.Logger.Fatalln(err)
	}
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
		log.Logger.Tracef("generating file path: %s", ltPath)
		dir := getBasedirectoryFromPath(ltPath)
		common.CreateDirectory(dir)
		if common.FileExists(ltPath) {
			log.Logger.Tracef("skipping %s because %s exists",
				name, ltPath)
			continue
		}
		_, err := GenLaunchTemplateFromInstanceId(region, id, ltPath)
		if err != nil {
			log.Logger.Fatalln(err)
		}
	}
}
