package cmd

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"os"

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

func PrintMap(m map[string]interface{}, prefix string) {
	for key, value := range m {
		fmt.Printf("%s%s: ", prefix, key)
		switch value.(type) {
		case map[string]interface{}:
			fmt.Println()
			PrintMap(value.(map[string]interface{}), prefix+"  ")
		default:
			fmt.Printf("%v\n", value)
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

	log.Logger.Printf("Original size: %d bytes\n", len(str))
	log.Logger.Printf("Compressed size: %d bytes\n", len(compressedData))
	log.Logger.Printf("Base64 encoded size: %d bytes\n", len(encodedData))
	log.Logger.Printf("Base64 encoded data: %s\n", encodedData)

	return encodedData, nil
}
