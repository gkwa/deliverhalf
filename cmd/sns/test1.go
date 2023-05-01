/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/spf13/cobra"
	meta "github.com/taylormonacelli/deliverhalf/cmd/meta"

	"github.com/spf13/viper"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test1",
	Short: "test message is fake data and varies only in epochtime",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		test1()
	},
}

func init() {
	snsCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func test1() {
	topicARN := viper.GetString("sns.topic-arn")
	topicRegion := viper.GetString("sns.region")

	jsonStr := meta.GenTestBlob()

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
