//lint:file-ignore U1000 Return to this when i've pulled my head out of my ass
package cmd

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details

	"gorm.io/gorm"

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
		GetIdentityDocFromSNS(region)
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

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Auto Migrate
	db.AutoMigrate(&myimds.IdentityBlob{})

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
			jsonStr, err := json.Marshal(message)
			if err != nil {
				log.Logger.Fatalf("failed to marshal message, error: %s", err)
			}
			mydb.WriteToDb(db, string(jsonStr))

			deleteMessage(message, sqsClient, &config)
		}
	}
}

func base64EncodeMessage(message types.Message) (string, error) {
	jsonStr, err := json.Marshal(message)
	if err != nil {
		log.Logger.Fatalf("Error serializing message: %s", err)
	}
	base64Str := base64.StdEncoding.EncodeToString([]byte(jsonStr))
	log.Logger.Tracef("base64 encoded json message: %s", base64Str)
	return string(base64Str), err
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
