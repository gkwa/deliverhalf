/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// fakeCmd represents the fake command
var fakeCmd = &cobra.Command{
	Use:   "fake",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// logger := setupLogger()
		fmt.Println(getTestBlobBase64())
	},
}

func init() {
	metaCmd.AddCommand(fakeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fakeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fakeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getTestBlobBase64() string {
	result := GenTestBlob()
	msg := []byte(result)
	base64Str := base64.StdEncoding.EncodeToString(msg)
	return base64Str
}

func GenTestBlob() string {
	jsonStr := `{
        "accountId": "348759328109",
        "architecture": "arm64",
        "availabilityZone": "us-east-1c",
        "billingProducts": [
            "bp-8f5a09f1"
        ],
        "devpayProductCodes": null,
        "epochtime": %d,
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

	// Get the current Unix Epoch time
	epoch := time.Now().Unix()

	// Format the JSON string with the epoch time
	formattedJson := fmt.Sprintf(jsonStr, epoch)
	return string(formattedJson)
}
