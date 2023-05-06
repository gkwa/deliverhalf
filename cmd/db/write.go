/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/k0kubun/pp"
	"github.com/spf13/cobra"
	common "github.com/taylormonacelli/deliverhalf/cmd/common"
	imds "github.com/taylormonacelli/deliverhalf/cmd/ec2/imds"
	log "github.com/taylormonacelli/deliverhalf/cmd/logging"
	"gorm.io/gorm"

	// "gorm.io/driver/sqlite" // Sqlite driver based on GGO
	"github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
)

// writeCmd represents the write command
var writeCmd = &cobra.Command{
	Use:   "write",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger.Println("config called")

		fmt.Println("write called")
		// test(logger)
		Test2()
	},
}

func init() {
	dbCmd.AddCommand(writeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// writeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// writeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getSnsMessageFromStr(str string) (types.Message, error) {
	var message types.Message

	err := json.Unmarshal([]byte(str), &message)
	if err != nil {
		panic(fmt.Errorf("can't unmarshal %s into types.Message", str))
	}

	return message, nil
}

func Doit(db *gorm.DB, b64msg string) {
	compressed, err := common.CompresStrToB64(b64msg)
	if err != nil {
		log.Logger.Fatalf("cant compress %s: %s", b64msg, err)
	}

	// Decode the string
	decoded, err := base64.StdEncoding.DecodeString(b64msg)
	if err != nil {
		log.Logger.Fatal("Failed to decode base64 string")
	}

	message, err := getSnsMessageFromStr(string(decoded))
	if err != nil {
		log.Logger.Fatal(err)
	}

	body := make(map[string]interface{})
	err = json.Unmarshal([]byte(*message.Body), &body)
	if err != nil {
		log.Logger.Fatal(err)
	}

	subMessageBytes, err := base64.StdEncoding.DecodeString(body["Message"].(string))
	if err != nil {
		log.Logger.Fatal("Failed to decode base64 string")
	}
	subMessage := string(subMessageBytes)

	var doc imds.ExtendedInstanceIdentityDocument
	err = json.Unmarshal([]byte(subMessage), &doc)
	if err != nil {
		log.Logger.Fatal(err)
	}

	result := db.Create(&IdentityBlob{
		Doc:                     doc,
		B64SNSMessage:           b64msg,
		B64SNSMessageCompressed: compressed,
	})
	if result.Error != nil {
		log.Logger.Println(result.Error)
	}

	log.Logger.Printf("%d rows affected", result.RowsAffected)

	var blobs []IdentityBlob
	db.Find(&blobs)
	pp.Println(blobs)

	jsonBytes, err := json.MarshalIndent(blobs, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
	}
	fmt.Println(string(jsonBytes))
}

func Test2() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Auto Migrate
	db.AutoMigrate(&IdentityBlob{})

	// base64 encoded types.Message
	m1 := "eyJBdHRyaWJ1dGVzIjpudWxsLCJCb2R5Ijoie1xuICBcIlR5cGVcIiA6IFwiTm90aWZpY2F0aW9uXCIsXG4gIFwiTWVzc2FnZUlkXCIgOiBcImQ3ZGE3ZTFkLWY3MGQtNTI2MC1hMWJkLTIyZDNjMDVjODhlYVwiLFxuICBcIlRvcGljQXJuXCIgOiBcImFybjphd3M6c25zOnVzLXdlc3QtMjoxOTMwNDg4OTU3Mzc6ZGVsaXZlcmhhbGZcIixcbiAgXCJNZXNzYWdlXCIgOiBcImV3b2dJQ0FnSUNBZ0lDSmhZMk52ZFc1MFNXUWlPaUFpTXpRNE56VTVNekk0TVRBNUlpd0tJQ0FnSUNBZ0lDQWlZWEpqYUdsMFpXTjBkWEpsSWpvZ0ltRnliVFkwSWl3S0lDQWdJQ0FnSUNBaVlYWmhhV3hoWW1sc2FYUjVXbTl1WlNJNklDSjFjeTFsWVhOMExURmpJaXdLSUNBZ0lDQWdJQ0FpWW1sc2JHbHVaMUJ5YjJSMVkzUnpJam9nV3dvZ0lDQWdJQ0FnSUNBZ0lDQWlZbkF0T0dZMVlUQTVaakVpQ2lBZ0lDQWdJQ0FnWFN3S0lDQWdJQ0FnSUNBaVpHVjJjR0Y1VUhKdlpIVmpkRU52WkdWeklqb2diblZzYkN3S0lDQWdJQ0FnSUNBaVpYQnZZMmgwYVcxbElqb2dNVFk0TXpJME1ESXdPU3dLSUNBZ0lDQWdJQ0FpYVcxaFoyVkpaQ0k2SUNKaGJXa3RNR1kwT0RNMlpUQTVNRGxtTnpNeE5XWWlMQW9nSUNBZ0lDQWdJQ0pwYm5OMFlXNWpaVWxrSWpvZ0lta3RNRE00T0RnME4yUm1abVUxT0dSaE5ESWlMQW9nSUNBZ0lDQWdJQ0pwYm5OMFlXNWpaVlI1Y0dVaU9pQWliVFZoTGpSNGJHRnlaMlVpTEFvZ0lDQWdJQ0FnSUNKclpYSnVaV3hKWkNJNklHNTFiR3dzQ2lBZ0lDQWdJQ0FnSW0xaGNtdGxkSEJzWVdObFVISnZaSFZqZEVOdlpHVnpJam9nYm5Wc2JDd0tJQ0FnSUNBZ0lDQWljR1Z1WkdsdVoxUnBiV1VpT2lBaU1qQXlNeTB3TkMweU9WUXhOVG8wTlRveU0xb2lMQW9nSUNBZ0lDQWdJQ0p3Y21sMllYUmxTWEFpT2lBaU1UQXVNUzR5TGpFMUlpd0tJQ0FnSUNBZ0lDQWljbUZ0WkdsemEwbGtJam9nYm5Wc2JDd0tJQ0FnSUNBZ0lDQWljbVZuYVc5dUlqb2dJblZ6TFdWaGMzUXRNU0lzQ2lBZ0lDQWdJQ0FnSW5abGNuTnBiMjRpT2lBaU1qQXlNaTB4TVMwd055SUtJQ0FnSUgwPVwiLFxuICBcIlRpbWVzdGFtcFwiIDogXCIyMDIzLTA1LTA0VDIyOjQzOjI5LjM3OVpcIixcbiAgXCJTaWduYXR1cmVWZXJzaW9uXCIgOiBcIjFcIixcbiAgXCJTaWduYXR1cmVcIiA6IFwiRGp1OFY4SUZnM1FMZHZNeGxocThwWFBPZ0hRVUd3aG1HQjRqWXRqZjF1Z0JIeVR4TVdXcEFRdDJmVlFPQzNwL01yMUs0c3k3N1NuWmtPNWpHUlpyRC93WERDNVZ1SWY0QlA1bFZUZDdhNlNvWWhqcCtGMGRhNC9BRHVwdWtXWEpvSSt1QmhEWGtoVC9UMFZYRHJ2bWpXK0hLQUE3MVBISk04U0lWTXJYK1p3MU5Pd2ZsNTErTEt4N3hIN0QvWUNFR3RnalNOdUoyd0NaWVBQZUNPNlNEZGRRMFNDT0ZZcU9SV1l1TnBKQyswc3IvL2FJVHg3VzRvN1h0dkRVc0FRL2xhTXdyMnlsNVJkMXRsTndxN01NTmN1WHMrN28rT2F3TUJhK29WaHNWRE1DTjFHTG9OMklqS2pjaXpTYUVvRU16OWluWlVKdldXazZLd3NudkhYcmNnPT1cIixcbiAgXCJTaWduaW5nQ2VydFVSTFwiIDogXCJodHRwczovL3Nucy51cy13ZXN0LTIuYW1hem9uYXdzLmNvbS9TaW1wbGVOb3RpZmljYXRpb25TZXJ2aWNlLTU2ZTY3ZmNiNDFmNmZlYzA5YjAxOTY2OTI2MjVkMzg1LnBlbVwiLFxuICBcIlVuc3Vic2NyaWJlVVJMXCIgOiBcImh0dHBzOi8vc25zLnVzLXdlc3QtMi5hbWF6b25hd3MuY29tLz9BY3Rpb249VW5zdWJzY3JpYmVcdTAwMjZTdWJzY3JpcHRpb25Bcm49YXJuOmF3czpzbnM6dXMtd2VzdC0yOjE5MzA0ODg5NTczNzpkZWxpdmVyaGFsZjoyMGQ2MWZhZC02Nzg2LTRlYzYtOGFmMS04NTBmMjYzZTdiYTRcIlxufSIsIk1ENU9mQm9keSI6IjIwOTQ3ZDFmNjMwOWM1NTNkNzU1NmE5MDk0Nzc1MTE0IiwiTUQ1T2ZNZXNzYWdlQXR0cmlidXRlcyI6bnVsbCwiTWVzc2FnZUF0dHJpYnV0ZXMiOm51bGwsIk1lc3NhZ2VJZCI6IjYxM2U4ZTA0LWNiNjctNGRmYy1iYTVhLWE2MmM2OWI4MjcxOCIsIlJlY2VpcHRIYW5kbGUiOiJBUUVCMW9pRUhoVVRjNkFqSDFIUEhSZWdDZnhJOCtVUTJZZFhaY2FtSHM1MUhJY1JrdFBNU2szaXM1L3JmVTYvY2NweUFVbzh2Ui80b2R1Mm5mNnhTK0plVFNRbEVXQnZWMzRYcjdLQ3JqU252dU9pQTFsWlJwaXRnK1hyMnZYWldhSWdaa2dPYnE4UjMyQ01udFBpQkRQU0VNMWFzUUxKVVMraVBrdVppZTFuYlpINnRlUjVNSDZoT2xWTjFuS1VSUm5NRnF5MDFqMVdYdlRWdEdZdEx6VEd4Ylk5QmJ3WnpMb1NiT1BIOTVGdkRGWWNyNEdyOUhQV0dQMU8xUWM4T1BNR09HODdIWVJuY3BVVmRGYXhXeGhwRzRJOVhXSXJUK0hlOE5jZmdCd2M0b24ySlFXNGpjeVAvUzlFd2UwTGI3c3ZWSmhtc0NxNm9zRXZEWVpVSkllK1IwMVJnQTA2dzB6WjNtVHh4SUJlYWdTS2o1akpIb3cxazhRSThMVzlScVJDd04rZDZGaW4vWjdTd2xHUlFvNHhwUT09In0="

	Doit(db, m1)
}
