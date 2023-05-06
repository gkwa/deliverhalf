/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/taylormonacelli/deliverhalf/cmd/logging"

	"github.com/spf13/cobra"
)

// clearCmd represents the clear command
var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("clear called")
		main()
	},
}

func init() {
	snsCmd.AddCommand(clearCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clearCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clearCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// SQSDeleteMessageAPI defines the interface for the GetQueueUrl and DeleteMessage functions.
// We use this interface to test the functions using a mocked service.
type SQSDeleteMessageAPI interface {
	GetQueueUrl(ctx context.Context,
		params *sqs.GetQueueUrlInput,
		optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)

	DeleteMessage(ctx context.Context,
		params *sqs.DeleteMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error)
}

// GetQueueURL gets the URL of an Amazon SQS queue.
// Inputs:
//
//	c is the context of the method call, which includes the AWS Region.
//	api is the interface that defines the method call.
//	input defines the input arguments to the service call.
//
// Output:
//
//	If success, a GetQueueUrlOutput object containing the result of the service call and nil.
//	Otherwise, nil and an error from the call to GetQueueUrl.
func GetQueueURL(c context.Context, api SQSDeleteMessageAPI, input *sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
	return api.GetQueueUrl(c, input)
}

// RemoveMessage deletes a message from an Amazon SQS queue.
// Inputs:
//
//	c is the context of the method call, which includes the AWS Region.
//	api is the interface that defines the method call.
//	input defines the input arguments to the service call.
//
// Output:
//
//	If success, a DeleteMessageOutput object containing the result of the service call and nil.
//	Otherwise, nil and an error from the call to DeleteMessage.
func RemoveMessage(c context.Context, api SQSDeleteMessageAPI, input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	return api.DeleteMessage(c, input)
}

func main() {
	queue := flag.String("q", "", "The name of the queue")
	messageHandle := flag.String("m", "", "The receipt handle of the message")
	flag.Parse()

	if *queue == "" || *messageHandle == "" {
		fmt.Println("You must supply a queue name (-q QUEUE) and message receipt handle (-m MESSAGE-HANDLE)")
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logging.Logger.Fatalf("configuration error %s", err)
	}

	client := sqs.NewFromConfig(cfg)

	qUInput := &sqs.GetQueueUrlInput{
		QueueName: queue,
	}

	// Get URL of queue
	result, err := GetQueueURL(context.TODO(), client, qUInput)
	if err != nil {
		fmt.Println("Got an error getting the queue URL:")
		fmt.Println(err)
		return
	}

	queueURL := result.QueueUrl

	dMInput := &sqs.DeleteMessageInput{
		QueueUrl:      queueURL,
		ReceiptHandle: messageHandle,
	}

	_, err = RemoveMessage(context.TODO(), client, dMInput)
	if err != nil {
		fmt.Println("Got an error deleting the message:")
		fmt.Println(err)
		return
	}

	fmt.Println("Deleted message from queue with URL " + *queueURL)
}
