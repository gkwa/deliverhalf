/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	log "github.com/taylormonacelli/deliverhalf/cmd/logging"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger.Traceln("upload called")
		upload()
	},
}

func init() {
	logsCmd.AddCommand(uploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func uploadFileToS3(bucketName, key, filePath string) error {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		return fmt.Errorf("failed to load AWS configuration: %v", err)
	}

	// Create an S3 client
	client := s3.NewFromConfig(cfg)

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Create an S3 PutObjectInput
	input := &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &key,
		Body:   file,
	}

	// Upload the file to S3
	_, err = client.PutObject(context.TODO(), input)
	if err != nil {
		log.Logger.Errorf("failed to upload file to S3: %v", err)
		return err
	}

	return nil
}

func upload() {
	bucketName := viper.GetString("s3bucket.name")
	key := "logs/app.log"
	filePath := "logs/app.log"

	log.Logger.Tracef("bucketName: %s, key: %s, filePath: %s",
		bucketName, key, filePath,
	)

	err := uploadFileToS3(bucketName, key, filePath)
	if err != nil {
		log.Logger.Errorf("Failed to upload file to S3: %s", err)
		return
	}

	log.Logger.Trace("File uploaded successfully!")
}
