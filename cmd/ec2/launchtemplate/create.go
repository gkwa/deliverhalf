/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/taylormonacelli/deliverhalf/cmd/logging"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger.Trace("create called")
		create()
	},
}

func init() {
	LaunchtemplateCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getConfig() (map[string]interface{}, error) {
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	// Return a map of all the settings loaded by Viper
	return viper.AllSettings(), nil
}

func create() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Logger.Fatal("failed to load SDK configuration, " + err.Error())
	}

	// Create a new EC2 client
	ec2Client := ec2.NewFromConfig(cfg)

	// Call getConfig() to retrieve the config data
	config, err := getConfig()
	if err != nil {
		log.Logger.Fatal(err)
	}

	// Get the value from Viper and convert it to the custom type
	var instanceType types.InstanceType

	// Create a new source with a specific seed
	source := rand.NewSource(time.Now().UnixNano())

	// Create a new random number generator using the source
	rng := rand.New(source)

	tplNamePrefix := "deliverhalf"

	ltNameFromConfig := "test1"
	// Generate a random number between 0 and 999999
	randomNumber := rng.Intn(1000000)

	// Create a string with the random number appended
	ltName := fmt.Sprintf("%s-%s-%06d", tplNamePrefix, ltNameFromConfig, randomNumber)

	configIndex := fmt.Sprintf("%s.%s", "launch_templates", ltNameFromConfig)

	lt, ok := config["launch_templates"].(map[string]interface{})[ltNameFromConfig].(map[string]interface{})
	if !ok {
		log.Logger.Fatalf("could not index %s",
			fmt.Sprintf("%s.%s from %s", "launch_templates", ltNameFromConfig, viper.ConfigFileUsed()))
	}

	err = viper.UnmarshalKey(lt["instancetype"].(string), &instanceType)
	if err != nil {
		log.Logger.Fatalf("failed to unmarshal template %s", ltNameFromConfig)
	}

	// Specify the details of the launch template to create
	imageID := lt["imageid"].(string)
	keyName := "my-key-pair"

	sgIndex := fmt.Sprintf("%s.securitygroupids", configIndex)

	securityGroupIDs := viper.GetStringSlice(sgIndex)
	if len(securityGroupIDs) == 0 {
		log.Logger.Warnf("no security groups matched %s from %s, using new security group",
			sgIndex, viper.ConfigFileUsed())
	}

	log.Logger.Trace("security groups: %s", strings.Join(securityGroupIDs, ", "))

	script, err := createUserData()
	if err != nil {
		log.Logger.Fatalf("can't create userdata script, error: %s", err)
	}
	userDataEncoded := base64.StdEncoding.EncodeToString([]byte(script))

	// Create the launch template
	createLaunchTemplateInput := &ec2.CreateLaunchTemplateInput{
		LaunchTemplateName: &ltName,
		LaunchTemplateData: &types.RequestLaunchTemplateData{
			ImageId:          &imageID,
			InstanceType:     instanceType,
			KeyName:          &keyName,
			SecurityGroupIds: securityGroupIDs,
			UserData:         &userDataEncoded,
		},
	}

	createLaunchTemplateOutput, err := ec2Client.CreateLaunchTemplate(context.Background(), createLaunchTemplateInput)
	if err != nil {
		log.Logger.Fatal("failed to create launch template, " + err.Error())
	}

	log.Logger.Trace("Launch template created with ID:", *createLaunchTemplateOutput.LaunchTemplate.LaunchTemplateId)
}
