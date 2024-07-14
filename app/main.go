package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func handler(_ context.Context) (string, error) {
	start := time.Now()

	// Recupera variáveis de ambiente
	uri := os.Getenv("MONGODB_URI")
	parallel := os.Getenv("MONGODB_PARALLEL")
	awsRegion := os.Getenv("AWS_REGION")

	timestamp := time.Now().Format("20060102T150405")
	archivePath := fmt.Sprintf("/tmp/%s.dump.gz", timestamp)
	s3Bucket := "ufabc-next"
	s3Key := fmt.Sprintf("mongodb-next-backup/%s.dump.gz", timestamp)

	// Step 1: Construct mongodump command
	command := exec.Command("mongodump", "--uri", uri, "--numParallelCollections", parallel, "--archive="+archivePath, "--gzip")
	log.Printf("Running mongodump command: %v", command.Args)

	// Step 2: Execute mongodump command
	output, err := command.CombinedOutput()
	if err != nil {
		log.Printf("Backup failed: %v", err)
		log.Printf("Stderr: %s", output)
		return "Backup failed", err
	}
	log.Printf("Backup completed successfully")
	log.Printf("Stdout: %s", output)

	// Step 3: Upload to S3
	err = uploadToS3(s3Bucket, s3Key, archivePath, awsRegion)
	if err != nil {
		log.Printf("Failed to upload to S3: %v", err)
		return "Failed to upload to S3", err
	}
	log.Printf("Upload to S3 completed successfully")

	return fmt.Sprintf("Backup completed successfully and uploaded to S3 in %d ms", time.Since(start).Milliseconds()), nil
}

func uploadToS3(bucket, key, filePath, region string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %q, %v", filePath, err)
	}
	defer file.Close()

	uploader := s3.New(sess)
	_, err = uploader.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to S3: %v", err)
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
