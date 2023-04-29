/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"

	"encoding/base64"
)

// testCmd represents the test command
var test2Cmd = &cobra.Command{
	Use:   "test2",
	Short: "test message is fake data and is static",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		test2()
	},
}

func init() {
	snsCmd.AddCommand(test2Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func test2() {
	topicARN := viper.GetString("sns.topic-arn")
	topicRegion := viper.GetString("sns.region")

	jsonStr := `{
        "accountId": "348759328109",
        "architecture": "arm64",
        "availabilityZone": "us-east-1c",
        "billingProducts": [
            "bp-8f5a09f1"
        ],
        "devpayProductCodes": null,
        "epochtime": 1682801174,
        "imageId": "ami-0f4836e0909f7315f",
        "instanceId": "i-0388847dffe58da42",
        "instanceType": "m5a.4xlarge",
        "kernelId": null,
        "marketplaceProductCodes": null,
        "pendingTime": "2023-04-29T15:45:23Z",
        "privateIp": "10.1.2.15",
        "ramdiskId": null,
        "region": "us-east-1",
        "version": "2022-11-07"
    }`

	msg := []byte(jsonStr)
	base64Str := base64.StdEncoding.EncodeToString(msg)

	fmt.Printf("region: %s", topicRegion)

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(topicRegion))
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := sns.NewFromConfig(cfg)

	input := &sns.PublishInput{
		Message:  &base64Str,
		TopicArn: &topicARN,
	}

	result, err := PublishMessage(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got an error publishing the message:")
		fmt.Println(err)
		return
	}

	fmt.Println("Message ID: " + *result.MessageId)
}
