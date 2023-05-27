//lint:file-ignore U1000 Return to this when i've pulled my head out of my ass
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	mysns "github.com/taylormonacelli/deliverhalf/cmd/sns"
)

// unmarshal1Cmd represents the unmarshal command
var unmarshal1Cmd = &cobra.Command{
	Use:   "unmarshal1",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("unmarshal called")
		testUnmarshalingSns()
	},
}

func init() {
	testingCmd.AddCommand(unmarshal1Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// unmarshalCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// unmarshalCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func testUnmarshalingSns() {
	jsonBody := `{
		"Type": "Notification",
		"MessageId": "bc38c0ff-7eb1-4412-8fe5-9de536c7b5b2",
		"TopicArn": "arn:aws:sns:us-west-2:012345123489:myword",
		"TopicArn": "arn:aws:sns:us-west-2:012345123489:myword",
		"Message": "{\"version\":\"0\",\"id\":\"9a62dc04-5e95-4e78-a3f7-8f4b4d0188e7\",\"detail-type\":\"EC2 Instance State-change Notification\",\"source\":\"aws.ec2\",\"account\":\"012345123489\",\"time\":\"2023-05-27T21:09:25Z\",\"region\":\"us-west-2\",\"resources\":[\"arn:aws:ec2:us-west-2:012345123489:instance/i-8a72e6d6b642efc0b\"],\"detail\":{\"instance-id\":\"i-8a72e6d6b642efc0b\",\"state\":\"running\"}}",
		"Timestamp": "2023-05-27T21:09:26.165Z",
		"SignatureVersion": "1",
		"Signature": "HRMYb8HRAatNT/zf5VM1d9h/wBAqe6ouRpHPPECX8nftwJXLDuDo2RZzK2IOACjEObBP3I4gvhlqvpqQVqE2YW6saHXhUkMcfjzh5GPMHv+NenDaM9lJ46jaHfHGcHbulyU/NRUwiLPQ+uADi1PFSVzpetKyRO/6h0kv6F0d5EVTWtF1JGv7FcdKG9IbuDLPz6hlQrwYtuxeyEn3zBMecWJbqEVTYSKFnGT+6/NWBWySYsFkyoxrAQY4X2a5sWbAFAgXODbxEoUdLWO5kQ8cbZ41YgRaHbeEr1OvTbKwq4qsDJ6PBoPbJsv64ypLyGfwJe/ojWGqBSu7+CP+XapxqQ==",
		"SigningCertURL": "https://sns.us-west-2.amazonaws.com/SimpleNotificationService-01d088a6f77103d0fe307c0069e40ed6.pem",
		"UnsubscribeURL": "https://sns.us-west-2.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-west-2:012345123489:myword:20d61fad-6786-4ec6-8af1-850f263e7ba4"
	}`

	var notification mysns.NotificationMessage
	err := json.Unmarshal([]byte(jsonBody), &notification)
	if err != nil {
		fmt.Println("Error deserializing JSON:", err)
		return
	}

	// Deserialize the nested "Message" field
	err = json.Unmarshal([]byte(notification.Message), &notification.MessageDetail)
	if err != nil {
		fmt.Println("Error deserializing Message field:", err)
		return
	}

	// Access the deserialized fields
	fmt.Println("Type:", notification.Type)
	fmt.Println("MessageId:", notification.MessageID)
	fmt.Println("TopicArn:", notification.TopicArn)
	fmt.Println("Timestamp:", notification.Timestamp)
	fmt.Println("SignatureVersion:", notification.SignatureVersion)
	fmt.Println("Signature:", notification.Signature)
	fmt.Println("SigningCertURL:", notification.SigningCertURL)
	fmt.Println("UnsubscribeURL:", notification.UnsubscribeURL)
	fmt.Println("MessageDetail - Version:", notification.MessageDetail.Version)
	fmt.Println("MessageDetail - ID:", notification.MessageDetail.ID)
	fmt.Println("MessageDetail - DetailType:", notification.MessageDetail.DetailType)
	fmt.Println("MessageDetail - Source:", notification.MessageDetail.Source)
	fmt.Println("MessageDetail - Account:", notification.MessageDetail.Account)
	fmt.Println("MessageDetail - Time:", notification.MessageDetail.Time)
	fmt.Println("MessageDetail - Region:", notification.MessageDetail.Region)
	fmt.Println("MessageDetail - Resources:", notification.MessageDetail.Resources)
	fmt.Println("MessageDetail - Detail - Instance ID:", notification.MessageDetail.Detail.InstanceID)
	fmt.Println("MessageDetail - Detail - State:", notification.MessageDetail.Detail.State)
}
