package cmd

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	awsv1 "github.com/aws/aws-sdk-go/aws"
	"github.com/spf13/cobra"
	"github.com/taylormonacelli/deliverhalf/cmd"
	log "github.com/taylormonacelli/deliverhalf/cmd/logging"
)

// ec2Cmd represents the ec2 command
var Ec2Cmd = &cobra.Command{
	Use:   "ec2",
	Args:  cobra.OnlyValidArgs,
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger.Traceln("ec2 called")
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
	cmd.RootCmd.AddCommand(Ec2Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ec2Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ec2Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func CreateConfig(region string) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return aws.Config{}, err
	}
	return cfg, err
}

func GetEc2Client(region string) (*ec2.Client, error) {
	config, err := CreateConfig(region)
	if err != nil {
		return nil, err
	}
	// Create an EC2 client
	return ec2.NewFromConfig(config), nil
}

// GetTagValue returns the value of the tag with the specified key.
// If no tag with the given key is found, it returns an empty string.
func GetTagValue(tags *[]types.Tag, tagKey string) string {
	for _, tag := range *tags {
		if awsv1.StringValue(tag.Key) == tagKey {
			return awsv1.StringValue(tag.Value)
		}
	}
	return ""
}

func GetTagSpecificationValue(tagSpecs *[]types.LaunchTemplateTagSpecification, tagKey string) string {
	for _, spec := range *tagSpecs {
		if awsv1.StringValue((*string)(&spec.ResourceType)) != "instance" {
			continue
		}
		for _, tag := range spec.Tags {
			if awsv1.StringValue(tag.Key) == tagKey {
				return awsv1.StringValue(tag.Value)
			}
		}
	}

	return ""
}
