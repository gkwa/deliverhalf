/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/k0kubun/pp"

	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
	common "github.com/taylormonacelli/deliverhalf/cmd/common"
	myec2 "github.com/taylormonacelli/deliverhalf/cmd/ec2"
)

// subscribeCmd represents the subscribe command
var subscribeCmd = &cobra.Command{
	Use:   "subscribe",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := common.SetupLogger()
		region := "us-west-2"
		getIdentityDocFromSNS(logger, region)
	},
}

func init() {
	snsCmd.AddCommand(subscribeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// subscribeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// subscribeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type snsSqsConfig struct {
	topicRegion string
	topicARN    string
	sqsQueueARN string
	sqsQueueURL string
}

func getIdentityDocFromSNS(logger *log.Logger, region string) (imds.InstanceIdentityDocument, error) {
	// Load the AWS configuration

	config := snsSqsConfig{
		topicRegion: viper.GetString("sns.region"),
		topicARN:    viper.GetString("sns.topic-arn"),
		sqsQueueARN: viper.GetString("sqs.queue-arn"),
		sqsQueueURL: viper.GetString("sqs.queue-url"),
	}

	cfg, err := myec2.CreateConfig(logger, config.topicRegion)
	if err != nil {
		logger.Fatalf("Could not create config %s", err)
	}

	snsClient := sns.NewFromConfig(cfg)

	subscribeOutput, err := snsClient.Subscribe(context.Background(), &sns.SubscribeInput{
		Protocol: aws.String("sqs"), // Use SQS as the protocol
		TopicArn: aws.String(config.topicARN),
		Endpoint: aws.String(config.sqsQueueARN), // Specify the ARN of your SQS queue
	})
	if err != nil {
		logger.Fatalf("failed to subscribe to SNS topic: %v", err)
	}

	// Print the subscription ARN
	logger.Printf("Subscribed to SNS topic with ARN %s\n", *subscribeOutput.SubscriptionArn)

	// Create an SQS client
	sqsClient := sqs.NewFromConfig(cfg)

	// Receive messages from the SQS queue
	for {
		receiveOutput, err := sqsClient.ReceiveMessage(context.Background(), &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(config.sqsQueueURL), // Specify the URL of your SQS queue
			MaxNumberOfMessages: 10,                             // Change this to the desired number of messages to receive
			WaitTimeSeconds:     20,                             // Change this to the desired wait time
		})
		if err != nil {
			logger.Fatalf("failed to receive SQS message: %v", err)
		}

		for _, message := range receiveOutput.Messages {
			pp.Printf("Received message: %s", message)
			logger.Printf("Received message body: %s\n", *message.Body)
			doc, err := getIdentityDoc(logger, message)
			if err != nil {
				logger.Fatalf("failed trying to get identity document: %s", err)
			}
			pp.Println(doc)
			deleteMessage(logger, message, sqsClient, &config)
		}
	}
}

func deleteMessage(logger *log.Logger, message types.Message, client *sqs.Client, config *snsSqsConfig) {
	// Delete the message from the queue
	_, err := client.DeleteMessage(context.Background(), &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(config.sqsQueueURL), // Specify the URL of your SQS queue
		ReceiptHandle: message.ReceiptHandle,          // Use the receipt handle to identify the message
	})
	if err != nil {
		logger.Fatalf("failed to delete SQS message: %v", err)
	}
}

func getIdentityDoc(logger *log.Logger, message types.Message) (imds.InstanceIdentityDocument, error) {
	myStr := *message.Body
	x2, err := genMessageFromStr(logger, myStr)
	if err != nil {
		logger.Fatalf("error unmarshalling types.Message from %s: %s", myStr, err)
	}
	b64Str := *x2.Message

	decoded, err := b64DecodeStr(logger, b64Str)
	if err != nil {
		logger.Fatalf("could not base64 decode %s: %s", *message.Body, err)
	}

	doc, err := genIdentityDocFromJsonStr(logger, decoded)
	if err != nil {
		logger.Fatalf("could not generate Identiydocument from %s: %s", decoded, err)
	}
	return doc, err
}

func genMessageFromStr(logger *log.Logger, str string) (sns.PublishInput, error) {
	var pi sns.PublishInput

	err := json.Unmarshal([]byte(str), &pi)
	if err != nil {
		logger.Fatalf("unmarshalling %s into an sns.PublishInput failed: %s", str, err)
	}
	return pi, err
}

func b64DecodeStr(logger *log.Logger, b64 string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		logger.Fatalf("Failed to decode base64 string '%s': %s", b64, err)
	}
	return string(decoded), err
}

func genIdentityDocFromJsonStr(logger *log.Logger, jsStr string) (imds.InstanceIdentityDocument, error) {
	var doc imds.InstanceIdentityDocument

	err := json.Unmarshal([]byte(jsStr), &doc)
	if err != nil {
		logger.Fatalf("unmarshalling %s into an imds.InstanceIdentityDocument failed: %s", jsStr, err)
	}
	return doc, err
}
