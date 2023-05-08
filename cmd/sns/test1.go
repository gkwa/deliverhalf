/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/base64"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/spf13/cobra"
	log "github.com/taylormonacelli/deliverhalf/cmd/logging"
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

	log.Logger.Tracef("region: %s", topicRegion)

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(topicRegion))
	if err != nil {
		log.Logger.Fatalf("configuration error %s", err)
	}

	client := sns.NewFromConfig(cfg)

	input := &sns.PublishInput{
		Message:  &base64Str,
		TopicArn: &topicARN,
	}

	result, err := PublishMessage(context.TODO(), client, input)
	if err != nil {
		log.Logger.Traceln("Got an error publishing the message:")
		log.Logger.Traceln(err)
		return
	}

	log.Logger.Traceln("Message ID: " + *result.MessageId)
}

func SendJsonStr(jsonStr string) {
	topicARN := viper.GetString("sns.topic-arn")
	topicRegion := viper.GetString("sns.region")

	msg := []byte(jsonStr)
	base64Str := base64.StdEncoding.EncodeToString(msg)

	log.Logger.Printf("region: %s", topicRegion)

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(topicRegion))
	if err != nil {
		log.Logger.Fatalf("failed to load config: %s", err)
	}

	client := sns.NewFromConfig(cfg)

	input := &sns.PublishInput{
		Message:  &base64Str,
		TopicArn: &topicARN,
	}

	result, err := PublishMessage(context.TODO(), client, input)
	if err != nil {
		log.Logger.Printf("Got an error publishing the message: %s", err)
		return
	}

	log.Logger.Printf("Message ID: %s", *result.MessageId)
}
