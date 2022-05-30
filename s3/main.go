package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func main() {

	ctx := context.Background()

	f, err := os.Open("test.txt")
	if err != nil {
		panic(err)
	}

	s3session := session.Must(session.NewSession(
		&aws.Config{
			Region: aws.String("ap-south-1"),
		},
	))

	uploader := s3manager.NewUploader(s3session)

	bucketName := "s3-testing-2k22"
	key := "folder2/" + "folder3/" + time.Now().String() + ".txt"

	result, err := uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   f,
	})

	if err != nil {
		panic(err)
	}

	log.Println("File Uploaded Successfully, URL : ", result.Location)

}
