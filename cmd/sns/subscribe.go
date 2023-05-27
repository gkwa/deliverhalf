//lint:file-ignore U1000 Return to this when i've pulled my head out of my ass
package cmd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	// Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details

	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
	mydb "github.com/taylormonacelli/deliverhalf/cmd/db"
	myec2 "github.com/taylormonacelli/deliverhalf/cmd/ec2"
	myimds "github.com/taylormonacelli/deliverhalf/cmd/ec2/imds"
	log "github.com/taylormonacelli/deliverhalf/cmd/logging"
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
		region := "us-west-2"
		SaveSnsMessage(region)
	},
}

func init() {
	SnsCmd.AddCommand(subscribeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// subscribeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// subscribeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	mydb.Db.AutoMigrate(&myimds.IdentityBlob{})
}

type snsSqsConfig struct {
	topicRegion string
	topicARN    string
	sqsQueueARN string
	sqsQueueURL string
}

func SaveSnsMessage(region string) error {
	// Load the AWS configuration
	log.Logger.Debug("debugging is fun")

	config := snsSqsConfig{
		topicRegion: viper.GetString("sns.region"),
		topicARN:    viper.GetString("sns.topic-arn"),
		sqsQueueARN: viper.GetString("sqs.queue-arn"),
		sqsQueueURL: viper.GetString("sqs.queue-url"),
	}

	cfg, err := myec2.CreateConfig(config.topicRegion)
	if err != nil {
		log.Logger.Fatalf("Could not create config %s", err)
	}

	snsClient := sns.NewFromConfig(cfg)

	inputSns := &sns.SubscribeInput{
		Protocol: aws.String("sqs"), // Use SQS as the protocol
		TopicArn: aws.String(config.topicARN),
		Endpoint: aws.String(config.sqsQueueARN), // Specify the ARN of your SQS queue
	}

	subscribeOutput, err := snsClient.Subscribe(context.Background(), inputSns)
	if err != nil {
		log.Logger.Fatalf("failed to subscribe to SNS topic: %v", err)
	}

	// Print the subscription ARN
	log.Logger.Debugf("Subscribed to SNS topic with ARN %s", *subscribeOutput.SubscriptionArn)

	// Create an SQS client
	sqsClient := sqs.NewFromConfig(cfg)

	inputSqs := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(config.sqsQueueURL),
		MaxNumberOfMessages: 10, // Change this to the desired number of messages to receive
		WaitTimeSeconds:     20, // Change this to the desired wait time
	}

	// Receive messages from the SQS queue
	for {
		receivedOut, err := sqsClient.ReceiveMessage(context.Background(), inputSqs)
		if err != nil {
			log.Logger.Fatalf("failed to receive SQS message: %v", err)
		}
		persistReceivedOutput(receivedOut)
		processReceivedOutput(receivedOut)
		popReceivedOutput(receivedOut, sqsClient, &config)
	}
}

func popReceivedOutput(receiveOut *sqs.ReceiveMessageOutput, sqsClient *sqs.Client, snsSqsConfig *snsSqsConfig) {
	for _, message := range receiveOut.Messages {
		// Delete the message from the queue
		_, err := sqsClient.DeleteMessage(context.Background(), &sqs.DeleteMessageInput{
			QueueUrl:      aws.String(snsSqsConfig.sqsQueueURL), // Specify the URL of your SQS queue
			ReceiptHandle: message.ReceiptHandle,                // Use the receipt handle to identify the message
		})
		if err != nil {
			log.Logger.Fatalf("failed to delete SQS message: %v", err)
		}
	}
}

func persistReceivedOutput(receiveOut *sqs.ReceiveMessageOutput) {
	jsonReceiveOut, err := json.MarshalIndent(receiveOut, "", "  ")
	if err != nil {
		log.Logger.Fatal(err)
	}

	x1 := ExtendedSqsReceiveMessageOutput{JsonDef: string(jsonReceiveOut)}
	result := mydb.Db.Create(&x1)
	if result.Error != nil {
		log.Logger.Fatalln("Error:", result.Error)
	}
}

func processReceivedOutput(receiveOutput *sqs.ReceiveMessageOutput) {
	out, err := json.MarshalIndent(receiveOutput, "", "  ")
	if err != nil {
		log.Logger.Fatal(err)
	}
	fmt.Print(string(out))

	for _, message := range receiveOutput.Messages {
		// *mesage.body is a NotificationMessage
		fmt.Printf("Received message body: %s", *message.Body)

		var x1 NotificationMessage
		json.Unmarshal([]byte(*message.Body), &x1)

		jsonData, err := json.MarshalIndent(x1, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling struct to JSON:", err)
			return
		}

		// Print the JSON data
		fmt.Printf("fartface: %s", string(jsonData))

		// -------------------

		var event myec2.EC2StateChangeEvent
		err = json.Unmarshal([]byte(x1.Message), &event)
		if err != nil {
			fmt.Println("Error deserializing message:", err)
			return
		}

		jsonData, err = json.MarshalIndent(event, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling struct to JSON:", err)
			return
		}

		// Print the JSON data
		fmt.Println(string(jsonData))

		// Access the deserialized fields
		fmt.Println("DetailType:", event.DetailType)

		// Check the detail-type value against the enum
		switch event.DetailType {
		case string(myec2.EC2StateChangeNotification):
			fmt.Println("Received EC2 State Change Notification")
		case string(myec2.OtherDetailType):
			fmt.Println("Received Other Detail Type")
		default:
			fmt.Println("Received Unknown Detail Type")
		}

	}
}

func GetIdentityDocFromSNS(region string) (imds.InstanceIdentityDocument, error) {
	// Load the AWS configuration
	log.Logger.Debug("debugging is fun")

	config := snsSqsConfig{
		topicRegion: viper.GetString("sns.region"),
		topicARN:    viper.GetString("sns.topic-arn"),
		sqsQueueARN: viper.GetString("sqs.queue-arn"),
		sqsQueueURL: viper.GetString("sqs.queue-url"),
	}

	cfg, err := myec2.CreateConfig(config.topicRegion)
	if err != nil {
		log.Logger.Fatalf("Could not create config %s", err)
	}

	snsClient := sns.NewFromConfig(cfg)

	subscribeOutput, err := snsClient.Subscribe(context.Background(), &sns.SubscribeInput{
		Protocol: aws.String("sqs"), // Use SQS as the protocol
		TopicArn: aws.String(config.topicARN),
		Endpoint: aws.String(config.sqsQueueARN), // Specify the ARN of your SQS queue
	})
	if err != nil {
		log.Logger.Fatalf("failed to subscribe to SNS topic: %v", err)
	}

	// Print the subscription ARN
	log.Logger.Debugf("Subscribed to SNS topic with ARN %s", *subscribeOutput.SubscriptionArn)

	// Create an SQS client
	sqsClient := sqs.NewFromConfig(cfg)

	// Auto Migrate

	// Receive messages from the SQS queue
	for {
		receiveOutput, err := sqsClient.ReceiveMessage(context.Background(), &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(config.sqsQueueURL), // Specify the URL of your SQS queue
			MaxNumberOfMessages: 10,                             // Change this to the desired number of messages to receive
			WaitTimeSeconds:     20,                             // Change this to the desired wait time
		})
		if err != nil {
			log.Logger.Fatalf("failed to receive SQS message: %v", err)
		}

		for _, message := range receiveOutput.Messages {
			log.Logger.Tracef("Received message body: %s", *message.Body)
			jsonStr, err := json.MarshalIndent(message, "", "  ")
			if err != nil {
				log.Logger.Fatalf("failed to marshal message, error: %s", err)
			}
			mydb.WriteToDb(mydb.Db, string(jsonStr))

			deleteMessage(message, sqsClient, &config)
		}
	}
}

func deleteMessage(message types.Message, client *sqs.Client, config *snsSqsConfig) {
	// Delete the message from the queue
	_, err := client.DeleteMessage(context.Background(), &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(config.sqsQueueURL), // Specify the URL of your SQS queue
		ReceiptHandle: message.ReceiptHandle,          // Use the receipt handle to identify the message
	})
	if err != nil {
		log.Logger.Fatalf("failed to delete SQS message: %v", err)
	}
}

func base64EncodeMessage(message types.Message) (string, error) {
	jsonStr, err := json.MarshalIndent(message, "", "  ")
	if err != nil {
		log.Logger.Fatalf("Error serializing message: %s", err)
	}
	base64Str := base64.StdEncoding.EncodeToString([]byte(jsonStr))
	log.Logger.Tracef("base64 encoded json message: %s", base64Str)
	return string(base64Str), err
}

func getIdentityDoc(message types.Message) (imds.InstanceIdentityDocument, error) {
	myStr := *message.Body
	x2, err := genMessageFromStr(myStr)
	if err != nil {
		log.Logger.Fatalf("error unmarshalling types.Message from %s: %s", myStr, err)
	}
	b64Str := *x2.Message

	decoded, err := b64DecodeStr(b64Str)
	if err != nil {
		log.Logger.Fatalf("could not base64 decode %s: %s", *message.Body, err)
	}

	doc, err := genIdentityDocFromJsonStr(decoded)
	if err != nil {
		log.Logger.Fatalf("could not generate Identiydocument from %s: %s", decoded, err)
	}
	return doc, err
}

func genMessageFromStr(str string) (sns.PublishInput, error) {
	var pi sns.PublishInput

	err := json.Unmarshal([]byte(str), &pi)
	if err != nil {
		log.Logger.Fatalf("unmarshalling %s into an sns.PublishInput failed: %s", str, err)
	}
	return pi, err
}

func b64DecodeStr(b64 string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		log.Logger.Fatalf("Failed to decode base64 string '%s': %s", b64, err)
	}
	return string(decoded), err
}

func genIdentityDocFromJsonStr(jsStr string) (imds.InstanceIdentityDocument, error) {
	var doc imds.InstanceIdentityDocument

	err := json.Unmarshal([]byte(jsStr), &doc)
	if err != nil {
		log.Logger.Fatalf("unmarshalling %s into an imds.InstanceIdentityDocument failed: %s", jsStr, err)
	}
	return doc, err
}
