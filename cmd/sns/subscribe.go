/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
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
		fmt.Println("subscribe called")
		logger := common.SetupLogger()
		region := "us-west-2"
		subscribe(logger, region)
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

func printMap(m map[string]interface{}, prefix string) {
	for key, value := range m {
		fmt.Printf("%s%s: ", prefix, key)
		switch value.(type) {
		case map[string]interface{}:
			fmt.Println()
			printMap(value.(map[string]interface{}), prefix+"  ")
		default:
			fmt.Printf("%v\n", value)
		}
	}
}

func subscribe(logger *log.Logger, region string) {
	// Load the AWS configuration

	topicRegion := viper.GetString("sns.region")
	topicARN := viper.GetString("sns.topic-arn")
	sqsQueueARN := viper.GetString("sqs.queue-arn")
	sqsQueueURL := viper.GetString("sqs.queue-url")

	// create a map to hold the decoded JSON data
	data := make(map[string]interface{})

	// print the data
	fmt.Println(data)

	cfg, err := myec2.CreateConfig(logger, topicRegion)
	if err != nil {
		logger.Fatalf("Could not create config %s", err)
	}
	// Create an SNS client
	snsClient := sns.NewFromConfig(cfg)

	// Subscribe to the SNS topic
	subscribeOutput, err := snsClient.Subscribe(context.Background(), &sns.SubscribeInput{
		Protocol: aws.String("sqs"), // Use SQS as the protocol
		TopicArn: aws.String(topicARN),
		Endpoint: aws.String(sqsQueueARN), // Specify the ARN of your SQS queue
	})
	if err != nil {
		panic(fmt.Sprintf("failed to subscribe to SNS topic: %v", err))
	}

	// Print the subscription ARN
	fmt.Printf("Subscribed to SNS topic with ARN %s\n", *subscribeOutput.SubscriptionArn)

	// Create an SQS client
	sqsClient := sqs.NewFromConfig(cfg)

	// Receive messages from the SQS queue
	for {
		receiveOutput, err := sqsClient.ReceiveMessage(context.Background(), &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(sqsQueueURL), // Specify the URL of your SQS queue
			MaxNumberOfMessages: 10,                      // Change this to the desired number of messages to receive
			WaitTimeSeconds:     20,                      // Change this to the desired wait time
		})
		if err != nil {
			logger.Fatal(fmt.Sprintf("failed to receive SQS message: %v", err))
		}

		// Print the message body
		for _, message := range receiveOutput.Messages {
			fmt.Printf("Received message: %s\n", *message.Body)

			// unmarshal the JSON data into the map
			err = json.Unmarshal([]byte(*message.Body), &data)
			if err != nil {
				logger.Fatalf("unmarshalling sns message body causes error %s", err)
			}

			// printMap(data, "")
			// fmt.Println()
			m := data["Message"].(string)
			// Decode the string
			decoded, err := base64.StdEncoding.DecodeString(m)
			if err != nil {
				logger.Fatalf("Failed to decode base64 string: %s", err)
			}
			jsStr := string(decoded)

			// Print the decoded string
			fmt.Println("Decoded string:", jsStr)

			// unmarshal the JSON data into the map
			err = json.Unmarshal([]byte(jsStr), &data)
			if err != nil {
				logger.Fatalf("unmarshalling sns message body causes error %s", err)
			}

			value, ok := data["epochtime"]
			if ok {
				epoch := int64(value.(float64))
				t := time.Unix(epoch, 0)
				fmt.Println(t.Weekday().String(), t.Format("January 02, 2006 15:04:05"))
			} else {
				fmt.Printf("Key '%s' does not exist in the map\n", "epochtime")
			}

			// Delete the message from the queue
			_, err = sqsClient.DeleteMessage(context.Background(), &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(sqsQueueURL), // Specify the URL of your SQS queue
				ReceiptHandle: message.ReceiptHandle,   // Use the receipt handle to identify the message
			})
			if err != nil {
				logger.Fatal(fmt.Sprintf("failed to delete SQS message: %v", err))
			}
		}
	}
}
