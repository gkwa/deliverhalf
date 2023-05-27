//lint:file-ignore U1000 Return to this when i've pulled my head out of my ass
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	myec2 "github.com/taylormonacelli/deliverhalf/cmd/ec2"
)

// unmarshal1Cmd represents the unmarshal command
var unmarshal2Cmd = &cobra.Command{
	Use:   "unmarshal2",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("unmarshal called")
		testUnmarshalinMessageBody()
	},
}

func init() {
	testingCmd.AddCommand(unmarshal2Cmd)
}

func testUnmarshalinMessageBody() {
	message := `{"version":"0","id":"0bc35f4a-94d2-4578-b6de-35f6e4af493c","detail-type":"EC2 Instance State-change Notification","source":"aws.ec2","account":"238904731545","time":"2023-05-27T15:30:21Z","region":"us-west-2","resources":["arn:aws:ec2:us-west-2:238904731545:instance/i-067e0a3f3a5b0d6be"],"detail":{"instance-id":"i-067e0a3f3a5b0d6be","state":"running"}}`
	var event myec2.EC2StateChangeEvent
	err := json.Unmarshal([]byte(message), &event)
	if err != nil {
		fmt.Println("Error deserializing message:", err)
		return
	}

	fmt.Println("Version:", event.Version)
	fmt.Println("ID:", event.ID)
	fmt.Println("DetailType:", event.DetailType)
	fmt.Println("Source:", event.Source)
	fmt.Println("Account:", event.Account)
	fmt.Println("Time:", event.Time)
	fmt.Println("Region:", event.Region)
	fmt.Println("Resources:", event.Resources)
	fmt.Println("Instance ID:", event.Detail.InstanceID)
	fmt.Println("State:", event.Detail.State)
}
