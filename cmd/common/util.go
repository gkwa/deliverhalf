package cmd

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	log "github.com/taylormonacelli/deliverhalf/cmd/logging"
)

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func IsFileEmpty(filePath string) (bool, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil // File doesn't exist, so it's not empty
		}
		return false, err // Error occurred while accessing file
	}

	return fileInfo.Size() == 0, nil // Return true if file size is 0, indicating an empty file
}

func PrintMap(m map[string]interface{}, prefix string) {
	for key, value := range m {
		log.Logger.Tracef("%s%s: ", prefix, key)
		switch v := value.(type) {
		case map[string]interface{}:
			log.Logger.Traceln()
			PrintMap(v, prefix+"  ")
		default:
			log.Logger.Tracef("%v\n", v)
		}
	}
}

func CompresStrToB64(str string) (string, error) {
	// Create a buffer to write the compressed data to
	var buf bytes.Buffer

	// Create a gzip writer that writes to the buffer
	gz := gzip.NewWriter(&buf)

	// Write the string to the gzip writer
	if _, err := gz.Write([]byte(str)); err != nil {
		log.Logger.Fatal(err)
	}

	// Close the gzip writer to flush any remaining data
	if err := gz.Close(); err != nil {
		log.Logger.Fatal(err)
	}

	// Get the compressed data as a byte slice
	compressedData := buf.Bytes()

	// Base64 encode the compressed data
	encodedData := base64.StdEncoding.EncodeToString(compressedData)

	log.Logger.Tracef("Original size: %d bytes\n", len(str))
	log.Logger.Tracef("Compressed size: %d bytes\n", len(compressedData))
	log.Logger.Tracef("Base64 encoded size: %d bytes\n", len(encodedData))
	log.Logger.Tracef("Base64 encoded data: %s\n", encodedData)

	return encodedData, nil
}

func DecodeBase64String(encoded string) string {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		log.Logger.Fatalf("Failed to decode base64 string")
	}
	return string(decoded)
}

// BucketBasics encapsulates the Amazon Simple Storage Service (Amazon S3) actions
// used in the examples.
// It contains S3Client, an Amazon S3 service client that is used to perform bucket
// and object actions.
type BucketBasics struct {
	S3Client *s3.Client
}

// DownloadFile gets an object from a bucket and stores it in a local file.
func (basics BucketBasics) DownloadFile(bucketName string, objectKey string, fileName string) error {
	result, err := basics.S3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		log.Logger.Errorf("Couldn't get object %v:%v. Here's why: %v\n", bucketName, objectKey, err)
		return err
	}
	defer result.Body.Close()
	file, err := os.Create(fileName)
	if err != nil {
		log.Logger.Errorf("Couldn't create file %v. Here's why: %v\n", fileName, err)
		return err
	}
	defer file.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Logger.Errorf("Couldn't read object body from %v. Here's why: %v\n", objectKey, err)
	}
	_, err = file.Write(body)
	return err
}

func DownloadFileFromS3(downloader *s3.Client, bucketName, filename, localTargetPath string) error {
	basics := BucketBasics{S3Client: downloader}
	return basics.DownloadFile(bucketName, filename, localTargetPath)
}

func HandleDownloadError(err error) {
	if err != nil {
		log.Logger.Errorf("failed to download object from S3: %v", err)
	}
}
