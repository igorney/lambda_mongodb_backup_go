package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func handler(_ context.Context) (string, error) {
	start := time.Now()
	uri := os.Getenv("MONGODB_URI")
	parallel := os.Getenv("MONGODB_PARALLEL")
	timestamp := time.Now().Format("20060102T150405")
	archivePath := fmt.Sprintf("/app/backups/%s.dump.gz", timestamp)
	s3Bucket := "ufabc-next"
	s3Key := fmt.Sprintf("mongodb-next-backup/%s.dump.gz", timestamp)

	// Step 1: Construct mongodump command
	command := exec.Command("mongodump", "--uri", uri, "--numParallelCollections", parallel, "--archive="+archivePath, "--gzip")
	log.Printf("Running mongodump command: %v", command.Args)
	step1 := time.Now()

	// Step 2: Execute mongodump command
	output, err := command.CombinedOutput()
	step2 := time.Now()
	if err != nil {
		log.Printf("Backup failed: %v", err)
		log.Printf("Stderr: %s", output)
		return "Backup failed", err
	}
	log.Printf("Backup completed successfully")
	log.Printf("Stdout: %s", output)

	// Step 3: Upload to S3
	err = uploadToS3(s3Bucket, s3Key, archivePath)
	step3 := time.Now()
	if err != nil {
		log.Printf("Failed to upload to S3: %v", err)
		return "Failed to upload to S3", err
	}
	log.Printf("Upload to S3 completed successfully")

	// Log times
	log.Printf("Time to construct command: %d ms", step1.Sub(start).Milliseconds())
	log.Printf("Time to execute mongodump: %d ms", step2.Sub(step1).Milliseconds())
	log.Printf("Time for upload to S3: %d ms", step3.Sub(step2).Milliseconds())
	log.Printf("Total time: %d ms", step3.Sub(start).Milliseconds())

	return "Backup completed successfully and uploaded to S3", nil
}

func uploadToS3(bucket, key, filePath string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			os.Getenv("AWS_SESSION_TOKEN"),
		),
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
	// Execute localmente
	result, err := handler(context.Background())
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result)
	}
}
