package common

import (
	"io/ioutil"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Loads body of a given S3 object.
func LoadS3Object (bucket, key string) ([]byte, error) {
	objectIn := s3.GetObjectInput {
		Bucket: aws.String(bucket),
		Key: aws.String(key),
	}
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
