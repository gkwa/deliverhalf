//lint:file-ignore U1000 Return to this when i've pulled my head out of my ass

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/spf13/cobra"
	mysns "github.com/taylormonacelli/deliverhalf/cmd/sns"
)

// unmarshal3Cmd represents the unmarshal3 command
var unmarshal3Cmd = &cobra.Command{
	Use:   "unmarshal3",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("unmarshal3 called")
		testUnmarshalingSqsReceiveMessageOutput()
	},
}

func init() {
	testingCmd.AddCommand(unmarshal3Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// unmarshal3Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// unmarshal3Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func testUnmarshalingSqsReceiveMessageOutput() {
	output := `
	{
		"Messages": [
		  {
			"Attributes": null,
			"Body": "{\n  \"Type\" : \"Notification\",\n  \"MessageId\" : \"f2d1d36c-332b-5a4f-90f3-8f80a700cb40\",\n  \"TopicArn\" : \"arn:aws:sns:us-west-2:987654321098:deliverhalf\",\n  \"Message\" : \"{\\\"s3Bucket\\\":\\\"ibm-bigcat\\\",\\\"s3ObjectKey\\\":[\\\"AWSLogs/987654321098/CloudTrail/us-west-2/2023/05/27/987654321098_CloudTrail_us-west-2_20230527T2325Z_4Ef2boJWDeex6r8s.json.gz\\\"]}\",\n  \"Timestamp\" : \"2023-05-27T23:27:39.734Z\",\n  \"SignatureVersion\" : \"1\",\n  \"Signature\" : \"JdNYR43Xm/dZ3XmvOdqOIhZYOi1dOMpOIQLoPHAypiLy1EYegnoyiwF2shA5QlIpW2Q/x2Rxp+FbguJGhK+L0paREhDrf/h1AbdiTCc5GhgYm8rEjVDCjQHs9cuzEyJ+az5Yqq0flQ9BiFQ9/JIjR+wXPXBPUKYf0aIlTXcKDBB1MQY+MRWWPBts4QdaeNyv+PhfQHka5oRU6idW2glcOJOgzLVLEtDUWYyZDMemMpN5y+ZHwJWoaCSwqRMQ8iyRKQlUf7CJiEpzuFcBssZ1Xyrjn1Zx73qtXOIzt1sI3m86izGZzkjhi5mVcWlr1LpOGU4PA2hg7oj+q1JLMnUpcg==\",\n  \"SigningCertURL\" : \"https://sns.us-west-2.amazonaws.com/SimpleNotificationService-01d088a6f77103d0fe307c0069e40ed6.pem\",\n  \"UnsubscribeURL\" : \"https://sns.us-west-2.amazonaws.com/?Action=Unsubscribe\u0026SubscriptionArn=arn:aws:sns:us-west-2:987654321098:deliverhalf:20d61fad-6786-4ec6-8af1-850f263e7ba4\"\n}",
			"MD5OfBody": "87d88ef0d114a81900553224687da0e4",
			"MD5OfMessageAttributes": null,
			"MessageAttributes": null,
			"MessageId": "aeb2d209-9b76-49ab-87c9-91c7de2f0bbf",
			"ReceiptHandle": "AQEBVGDaPaqYZdD4KJLhY3C5FBt2KdlrtOu/ZzGtALuyDmAh08YAeCamL/g5HEOu07h4KbligWxVyszeLJ+GRtTcQKbzmPDNo0tj7sc4S2z3z41i2U0yCXDJ0sLSAKxLt9rbz2iemICC4DNdurOFCj1lcfjK1cIz5gRkrEtIsTKsNeWibbcE0IXGT8TUlSyWxLS+EJ6JQ/P+OWmT39oBjrDWzjENQHo4cUWaF+OaPs9412l2/ho3uedAFWOzw9e0dQihbdrjM+0gUn74BgiJvxAehuCOKA9XrI1evh3TOyp372Qa3cY6skZf+JhqN+U3unqXR2g8r4v0v3jAHSGO97u3GBU4gRgJM/3UEQkLY6aRE2JFxFg+I3qA9r4BOQ7KBThYqV3iR1LqT1L9f1lDASvPCw=="
		  },
		  {
			"Attributes": null,
			"Body": "{\n  \"Type\" : \"Notification\",\n  \"MessageId\" : \"4a82f501-244e-523f-abe8-a9edf6ce7f57\",\n  \"TopicArn\" : \"arn:aws:sns:us-west-2:987654321098:deliverhalf\",\n  \"Message\" : \"{\\\"s3Bucket\\\":\\\"ibm-bigcat\\\",\\\"s3ObjectKey\\\":[\\\"AWSLogs/987654321098/CloudTrail/us-east-1/2023/05/28/987654321098_CloudTrail_us-east-1_20230528T0350Z_2Yb4XwljNGoirpZz.json.gz\\\"]}\",\n  \"Timestamp\" : \"2023-05-28T03:53:36.660Z\",\n  \"SignatureVersion\" : \"1\",\n  \"Signature\" : \"bHJ5vHYhaq5Ck7OcEisunHWqO7glR2mrWiiH/U/Ehs1NRfbkuoiROcUks3yKvNHzTDapt7NcXYE0s3bMgbUghjE+RXYPw+UCf9WZtq/UZZGhkYD8LpT6CaE4cWAj7vy8baKVaKCt3xAqOcDSh5OeSU93Llm9BQ19Rrnox6PH6x+3PlEad+wxBA6M3Q68FK0L/9OWHjhtRTvmd3QPCo/VwMd/l9qcxMyrd1JyBpYLyuvFqVnlm2Rabz/hYCsp9Up9tGNnP/cE8ooGe4srBzsnqc7N4sqjWJm1GC/guvEXloPeHwgHKKOPgfVxXTLCZI0tn9xV8a4ODQ3cDnHzf8cWiw==\",\n  \"SigningCertURL\" : \"https://sns.us-west-2.amazonaws.com/SimpleNotificationService-01d088a6f77103d0fe307c0069e40ed6.pem\",\n  \"UnsubscribeURL\" : \"https://sns.us-west-2.amazonaws.com/?Action=Unsubscribe\u0026SubscriptionArn=arn:aws:sns:us-west-2:987654321098:deliverhalf:20d61fad-6786-4ec6-8af1-850f263e7ba4\"\n}",
			"MD5OfBody": "7c1b9dfe181115c8b3ae0d197ee4a8ca",
			"MD5OfMessageAttributes": null,
			"MessageAttributes": null,
			"MessageId": "0ecc7e84-9796-43d2-9e07-7e831801f716",
			"ReceiptHandle": "AQEBSQzqiLErqL+ZKc1+RY6x8jc2S9finaj/+kFtm4EO/SRtR+5/VP831kypzAL86mC3RZNQ4G9/o7pU3Mpzbq7jBihhly8IOmC1lvFEsRX25SUpZNRQvlm6qKLu1CEIwtfBB2yMoo5URLzgD1vI9646CVj3Of+8szME33c3sTFj88GUMlEXxwk4dUXoosP5cj8i7y9bCYLK2BXyrolR/wteZI9ss9xlP9hyvZwItSc2/xAGA8Xzx7w1deRpuDAlan1nIf2lrvV7TdVXnl+vQaHQclvikOW/7Zd4e1n85p56xr3s2Gdc97N1sFSL/3wLm2XPi9ULWIp/soGaJMoH0jwGGFpxi8RguBSogJtrrzcxsfmVE8x4JtHsQO1hkZE1VEDCqIwv2iR3IeDxOxPo1ZTUJg=="
		  }
		]
	}`

	var rmo sqs.ReceiveMessageOutput
	err := json.Unmarshal([]byte(output), &rmo)
	if err != nil {
		fmt.Println("Error deserializing message:", err)
		return
	}

	for _, message := range rmo.Messages {

		jsonBytes, _ := json.MarshalIndent(message, "", "  ")
		fmt.Println(string(jsonBytes))
		jsonBytes, _ = json.MarshalIndent(message.Attributes, "", "  ")
		fmt.Println("Attributes:", string(jsonBytes))
		fmt.Println("Body:", *message.Body)

		var notification mysns.NotificationMessage
		err := json.Unmarshal([]byte(*message.Body), &notification)
		if err != nil {
			fmt.Println("Error deserializing JSON:", err)
			return
		}

		// jsonBytes, _ = json.MarshalIndent(notification, "", "  ")
		// fmt.Println("Notification:", string(jsonBytes))

		// var data map[string]json.RawMessage
		// err = json.Unmarshal([]byte(notification.Message), &data)
		// if err != nil {
		// 	fmt.Println("Error unmarshaling JSON:", err)
		// 	return
		// }
		// jsonBytes, _ = json.MarshalIndent(data, "", "  ")
		// fmt.Println("Message:", string(jsonBytes))

		// fmt.Println("MD5OfBody:", *message.MD5OfBody)

		// if message.MD5OfMessageAttributes != nil {
		// 	fmt.Println("MD5OfMessageAttributes:", *message.MD5OfMessageAttributes)
		// }

		// jsonBytes, _ = json.MarshalIndent(message.MessageAttributes, "", "  ")
		// fmt.Println("MessageAttributes:", string(jsonBytes))

		// fmt.Println("MessageId:", *message.MessageId)
		// fmt.Println("ReceiptHandle:", *message.ReceiptHandle)
	}
}
