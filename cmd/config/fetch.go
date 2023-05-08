/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/spf13/cobra"
	common "github.com/taylormonacelli/deliverhalf/cmd/common"
	myec2 "github.com/taylormonacelli/deliverhalf/cmd/ec2"
	log "github.com/taylormonacelli/deliverhalf/cmd/logging"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger.Traceln("fetch called")
		region := "us-west-2"
		// test( region)
		fetchConfigFromS3(region)
	},
}

// BucketBasics encapsulates the Amazon Simple Storage Service (Amazon S3) actions
// used in the examples.
// It contains S3Client, an Amazon S3 service client that is used to perform bucket
// and object actions.
type BucketBasics struct {
	S3Client *s3.Client
}

func test(region string) {
	fetchConfigFromS3(region)
	reloadConfig()
	showSettings()
}

func init() {
	configCmd.AddCommand(fetchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fetchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fetchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// DownloadFile gets an object from a bucket and stores it in a local file.
func (basics BucketBasics) DownloadFile(bucketName string, objectKey string, fileName string) error {
	result, err := basics.S3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		log.Logger.Printf("Couldn't get object %v:%v. Here's why: %v\n", bucketName, objectKey, err)
		return err
	}
	defer result.Body.Close()
	file, err := os.Create(fileName)
	if err != nil {
		log.Logger.Printf("Couldn't create file %v. Here's why: %v\n", fileName, err)
		return err
	}
	defer file.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Logger.Printf("Couldn't read object body from %v. Here's why: %v\n", objectKey, err)
	}
	_, err = file.Write(body)
	return err
}

func fetchConfigFromS3(region string) {
	// Load the AWS SDK configuration from the environment or shared config file
	cfg, err := myec2.CreateConfig(region)
	if err != nil {
		log.Logger.Fatalf("Could not create config %s", err)
	}

	// Create a new S3 downloader
	downloader := s3.NewFromConfig(cfg)

	// Decode the string
	decoded, err := base64.StdEncoding.DecodeString("c3RyZWFtYm94LWRlbGl2ZXJoYWxm")
	if err != nil {
		log.Logger.Fatalf("Failed to decode base64 string")
	}
	bucketName := string(decoded)

	decoded, err = base64.StdEncoding.DecodeString("LmRlbGl2ZXJoYWxmLnlhbWw=")
	if err != nil {
		log.Logger.Fatalf("Failed to decode base64 string")
	}
	config := string(decoded)

	s3path := fmt.Sprintf("s3://%s/%s", bucketName, config)

	// Print the decoded string
	log.Logger.Printf("Decoded s3path: %s", s3path)

	// Create a file to write the downloaded data to
	file, err := os.Create(config)
	if err != nil {
		log.Logger.Fatalf("error creating file %s: %s", config, err)
	}
	defer file.Close()

	dir := os.Getenv("HOME")
	filename := config
	localTargetPath := filepath.Join(dir, filename)

	// Create a new BucketBasics object with the S3 client
	basics := BucketBasics{S3Client: downloader}

	if common.FileExists(localTargetPath) {
		log.Logger.Printf("Overwriting %s with file from %s", localTargetPath, s3path)
	}

	// Download the object from S3 and store it in a local file
	err = basics.DownloadFile(bucketName, filename, localTargetPath)
	if err != nil {
		log.Logger.Fatalf("failed to download object from S3: %v", err)
	}

	// Print success message
	log.Logger.Printf("%s downloaded to %s", s3path, localTargetPath)
}
