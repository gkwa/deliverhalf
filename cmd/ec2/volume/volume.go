/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/spf13/cobra"
	myec2 "github.com/taylormonacelli/deliverhalf/cmd/ec2"
	imds "github.com/taylormonacelli/deliverhalf/cmd/ec2/imds"
	log "github.com/taylormonacelli/deliverhalf/cmd/logging"
	meta "github.com/taylormonacelli/deliverhalf/cmd/meta"
)

// volumeCmd represents the volume command
var VolumeCmd = &cobra.Command{
	Use:   "volume",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger.Traceln("volume called")
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
	myec2.Ec2Cmd.AddCommand(VolumeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// volumeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// volumeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type VolumeTag struct {
	Tags []types.Tag
	Size int32
}

func getVolumesFromInstanceIdentity(doc imds.ExtendedInstanceIdentityDocument, volumes *[]types.Volume) error {
	region := doc.Region
	instanceId := doc.InstanceId

	// Load the AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Logger.Fatal(err)
	}

	// Create a new EC2 client
	svc := ec2.NewFromConfig(cfg)

	err = getVolumesForInstance(context.Background(), svc, instanceId, volumes)
	if err != nil {
		log.Logger.Fatalf("Getting volumes for instance %s failed with error %s", instanceId, err)
	}

	volumeTags := extractVolumeTags(*volumes)
	printVolumeTags(volumeTags)
	return err
}

func testListVolumes() error {
	jsonStr := `
	{
		"accountId": "9876543210",
		"architecture": "x86",
		"availabilityZone": "us-east-1b",
		"billingProducts": [
		"bp-12345678"
		],
		"devpayProductCodes": null,
		"fetchTimestamp": 1736815649,
		"imageId": "ami-0987654321",
		"instanceId": "i-0987654321",
		"instanceType": "t2.micro",
		"kernelId": null,
		"marketplaceProductCodes": null,
		"pendingTime": "2023-05-07T14:30:00Z",
		"privateIp": "10.0.0.1",
		"ramdiskId": null,
		"region": "us-east-1",
		"version": "2021-05-01"
		}
		`

	jsonStr = `{
		"accountId": "193048895737",
		"architecture": "x86_64",
		"availabilityZone": "ap-southeast-2a",
		"billingProducts": [
			"bp-6ba54002"
		],
		"devpayProductCodes": null,
		"fetchTimestamp": 1683470368,
		"imageId": "ami-0b6125e77f55f0eff",
		"instanceId": "i-0488845dadd58da52",
		"instanceType": "t3a.2xlarge",
		"kernelId": null,
		"marketplaceProductCodes": null,
		"pendingTime": "2023-04-03T14:05:38Z",
		"privateIp": "172.31.18.139",
		"ramdiskId": null,
		"region": "ap-southeast-2",
		"version": "2017-09-30"
	}`
	doc, err := meta.GetIdentityDocFromStr(jsonStr)
	if err != nil {
		log.Logger.Fatalf("cant create %T: %s", doc, err)
	}
	var volumes []types.Volume
	err = getVolumesFromInstanceIdentity(doc, &volumes)
	if err != nil {
		log.Logger.Error(err)
		return err
	}

	// Marshal the Person struct to JSON with indentation
	jsonData, err := json.MarshalIndent(volumes, "", "  ")
	if err != nil {
		log.Logger.Traceln("Error marshaling struct to JSON:", err)
	}
	log.Logger.Trace(string(jsonData))
	log.Logger.Trace(string(jsonData))
	return nil
}

func extractVolumeTags(volumes []types.Volume) map[string]VolumeTag {
	volumeTags := make(map[string]VolumeTag)

	for _, volume := range volumes {
		volumeID := *volume.VolumeId
		tags := volume.Tags
		size := *volume.Size

		volumeTag := VolumeTag{
			Tags: tags,
			Size: size,
		}

		volumeTags[volumeID] = volumeTag
	}

	return volumeTags
}

func printVolumeTags(volumeTags map[string]VolumeTag) {
	for volumeID, volumeTag := range volumeTags {
		log.Logger.Tracef("Volume ID: %s\n", volumeID)
		log.Logger.Tracef("Size: %d\n", volumeTag.Size)
		log.Logger.Tracef("Tags:\n")
		for _, tag := range volumeTag.Tags {
			log.Logger.Tracef("  - %s: %s\n", *tag.Key, *tag.Value)
		}
		log.Logger.Traceln()
	}
}

func getVolumesForInstance(ctx context.Context, svc *ec2.Client, instanceID string, volumes *[]types.Volume) error {
	input := &ec2.DescribeVolumesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("attachment.instance-id"),
				Values: []string{instanceID},
			},
		},
	}
	resp, err := svc.DescribeVolumes(ctx, input)
	if err != nil {
		return err
	}

	*volumes = append(*volumes, resp.Volumes...)
	return nil
}
