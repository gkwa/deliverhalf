/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/spf13/cobra"
	myec2 "github.com/taylormonacelli/deliverhalf/cmd/ec2"
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
		test()
	},
}

type AMI struct {
	Name       string
	ImageID    string
	SnapshotID string
	Region     string
}

func (ami AMI) String() string {
	return fmt.Sprintf("Name: %s, ImageID: %s, Region: %s, SnapshotID: %s",
		ami.Name,
		ami.ImageID,
		ami.Region,
		ami.SnapshotID,
	)
}

func init() {
	amiCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func test() {
	amiName := fmt.Sprintf("my-image-%08d", time.Now().Unix())
	amiMeta := AMI{
		Name:       amiName,
		SnapshotID: "snap-01094c684c68ff680",
		Region:     "us-west-2",
	}

	rio, err := createAMIFromSnapshot(&amiMeta)
	if err != nil {
		log.Logger.Error(err)
		panic(err)
	}

	tagAmi(rio, &amiMeta)
	log.Logger.Printf("Created AMI with properties %s", amiMeta)
}

func tagAmi(ec2RegisterImageOutput *ec2.RegisterImageOutput, amiMeta *AMI) error {
	svc, err := myec2.GetEc2Client(amiMeta.Region)
	if err != nil {
		log.Logger.Error(err)
		return err
	}

	// Create tags for the registered image
	tagName := "Name"

	tagInput := &ec2.CreateTagsInput{
		Resources: []string{*ec2RegisterImageOutput.ImageId},
		Tags: []types.Tag{
			{
				Key:   &tagName,
				Value: &amiMeta.Name,
			},
		},
	}

	_, err = svc.CreateTags(context.Background(), tagInput)
	if err != nil {
		log.Logger.Fatalf("Error creating tags for image %s: %v", *ec2RegisterImageOutput.ImageId, err)
		return err
	}
	return nil
}

func createAMIFromSnapshot(ami *AMI) (*ec2.RegisterImageOutput, error) {
	svc, err := myec2.GetEc2Client(ami.Region)
	if err != nil {
		log.Logger.Error(err)
		return &ec2.RegisterImageOutput{}, err
	}

	// Call the RegisterImage function to register the AMI image
	input := &ec2.RegisterImageInput{
		Name:         aws.String(ami.Name),
		Architecture: types.ArchitectureValuesX8664,
		BlockDeviceMappings: []types.BlockDeviceMapping{
			{
				DeviceName: aws.String("/dev/sda1"),
				Ebs: &types.EbsBlockDevice{
					SnapshotId:          aws.String(ami.SnapshotID),
					VolumeSize:          aws.Int32(30),
					DeleteOnTermination: aws.Bool(true),
				},
			},
		},
		Description:    aws.String("created by deliverhalf"),
		RootDeviceName: aws.String("/dev/sda1"),
	}

	result, err := svc.RegisterImage(context.Background(), input)
	if err != nil {
		log.Logger.Fatalf("Error registering new AMI with snapshotID %s: %s", ami.SnapshotID, err)
	}
	ami.ImageID = *result.ImageId

	return result, err
}
