/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/base64"
	"errors"
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

func genRanName(prefix string) string {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	randomNumber := rng.Intn(1000000)
	ltName := fmt.Sprintf("%s-%06d", prefix, randomNumber)
	return ltName
}

func getLtFromName(name string) (map[string]interface{}, error) {
	// Call getConfig() to retrieve the config data
	config, err := getConfig()
	if err != nil {
		log.Logger.Fatal(err)
	}

	lt, ok := config["launch_templates"].(map[string]interface{})[name].(map[string]interface{})
	if !ok {
		msg := fmt.Sprintf("could not index %s.%s from %s",
			"launch_templates", name, viper.ConfigFileUsed())
		return nil, errors.New(msg)
	}
	return lt, nil
}

func create() {
	ltNameFromConfig := "test1"
	lt, err := getLtFromName(ltNameFromConfig)
	if err != nil {
		log.Logger.Fatalf("lookup template from '%s' failed", ltNameFromConfig)
	}

	region := lt["region"].(string)
	ltName := genRanName("deliverhalf")

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		log.Logger.Fatal("failed to load SDK configuration, " + err.Error())
	}

	// Create a new EC2 client
	ec2Client := ec2.NewFromConfig(cfg)

	// Get the value from Viper and convert it to the custom type
	var instanceType types.InstanceType

	configIndex := fmt.Sprintf("%s.%s", "launch_templates", ltNameFromConfig)
	err = viper.UnmarshalKey(lt["instancetype"].(string), &instanceType)
	if err != nil {
		log.Logger.Fatalf("failed to unmarshal template %s", ltNameFromConfig)
	}

	// Specify the details of the launch template to create
	imageID := lt["imageid"].(string)
	keyName := "my-key-pair"

	sgIndex := fmt.Sprintf("%s.securitygroupids", configIndex)
	iamInstanceProfileIndex := fmt.Sprintf("%s.iaminstanceprofile", configIndex)
	log.Logger.Tracef("instance profile: %s", iamInstanceProfileIndex)

	instanceProfileName := viper.GetString(iamInstanceProfileIndex)
	securityGroupIDs := viper.GetStringSlice(sgIndex)
	if len(securityGroupIDs) == 0 {
		log.Logger.Warnf("no security groups matched %s from %s, using new security group",
			sgIndex, viper.ConfigFileUsed())
	}

	log.Logger.Tracef("security groups: %s", strings.Join(securityGroupIDs, ", "))

	script, err := createUserData()
	if err != nil {
		log.Logger.Fatalf("can't create userdata script, error: %s", err)
	}
	userDataEncoded := base64.StdEncoding.EncodeToString([]byte(script))

	// Create the launch template
	createLaunchTemplateInput := &ec2.CreateLaunchTemplateInput{
		LaunchTemplateName: &ltName,
		LaunchTemplateData: &types.RequestLaunchTemplateData{
			ImageId:            &imageID,
			InstanceType:       instanceType,
			KeyName:            &keyName,
			SecurityGroupIds:   securityGroupIDs,
			UserData:           &userDataEncoded,
			IamInstanceProfile: &types.LaunchTemplateIamInstanceProfileSpecificationRequest{Name: &instanceProfileName},
			InstanceMarketOptions: &types.LaunchTemplateInstanceMarketOptionsRequest{
				MarketType: "spot",
			},
		},
	}

	createLaunchTemplateOutput, err := ec2Client.CreateLaunchTemplate(context.Background(), createLaunchTemplateInput)
	if err != nil {
		log.Logger.Fatal("failed to create launch template, " + err.Error())
	}

	log.Logger.Trace("Launch template created with ID:", *createLaunchTemplateOutput.LaunchTemplate.LaunchTemplateId)
}
