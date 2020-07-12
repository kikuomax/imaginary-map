package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Environment variable names.
const GEO_JSON_BUCKET_NAME = "GEO_JSON_BUCKET_NAME"
const ISLANDS_GEO_JSON_VERSION = "ISLANDS_GEO_JSON_VERSION"

// GeoJSON file name.
const GEO_JSON_FILE_NAME = "islands.json"

// Common form of a GetTileEvent.
type GetTileEvent struct {
	Zoom int `json:"zoom"`
	X int `json:"x"`
	Y int `json:"y"`
}

// Loads the GeoJSON for islands.
func LoadGeoJson () ([]byte, error) {
	bucketName, ok := os.LookupEnv(GEO_JSON_BUCKET_NAME)
	if !ok {
		return nil, errors.New(
			fmt.Sprintf(
				"environment variable %s is not set",
				GEO_JSON_BUCKET_NAME))
	}
	version, ok := os.LookupEnv(ISLANDS_GEO_JSON_VERSION)
	if !ok {
		return nil, errors.New(
			fmt.Sprintf(
				"environment variable %s is not set",
				ISLANDS_GEO_JSON_VERSION))
	}
	key := fmt.Sprintf("/%s/%s", version, GEO_JSON_FILE_NAME)
	objectIn := s3.GetObjectInput {
		Bucket: aws.String(bucketName),
		Key: aws.String(key),
	}
	log.Println("loading GeoJSON", bucketName, key)
	client := s3.New(session.New())
	objectOut, err := client.GetObject(&objectIn)
	if err != nil {
		return nil, err
	}
	bodyIn := objectOut.Body
	defer bodyIn.Close()
	body, err := ioutil.ReadAll(bodyIn)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func HandleRequest (ctx context.Context, event GetTileEvent) ([]byte, error) {
	return LoadGeoJson()
}

func main () {
	runtime.Start(HandleRequest)
}
