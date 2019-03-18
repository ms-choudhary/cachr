package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func keyExists(session *session.Session, bucket, key string) (bool, error) {
	svc := s3.New(session)
	results, err := svc.ListObjects(&s3.ListObjectsInput{Bucket: aws.String(bucket), Prefix: aws.String(key)})
	if err != nil {
		return false, err
	}

	return len(results.Contents) > 0, nil
}

func downloadS3(session *session.Session, bucket, key string) error {
	downloader := s3manager.NewDownloader(session)
	filename := filepath.Base(key)

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("cannot create file %s: %v", filename, err)
	}

	defer file.Close()

	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		return fmt.Errorf("cannot download key %s from s3: %v", key, err)
	}
	return nil
}

func getKey(session *session.Session, bucket, key string) error {
	if err := downloadS3(session, bucket, key); err != nil {
		return fmt.Errorf("cannot download from s3: %v", err)
	}

	if err := extractZipFile(filepath.Base(key), "."); err != nil {
		return fmt.Errorf("cannot extract zip file: %v", err)
	}

	return nil
}

func uploadS3(session *session.Session, bucket, key, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("cannot open file for saving: %v", err)
	}

	defer file.Close()

	uploader := s3manager.NewUploader(session)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("cannot upload file to s3: %v", err)
	}

	return nil
}

func saveKey(session *session.Session, bucket, key string, files []string) error {
	zipfilename := filepath.Base(key)

	if err := createZipFile(zipfilename, files); err != nil {
		return fmt.Errorf("cannot create zip file: %v", err)
	}

	if err := uploadS3(session, bucket, key, zipfilename); err != nil {
		return fmt.Errorf("cannot upload to s3: %v", err)
	}

	return nil
}

var CACHR_BUCKET = os.Getenv("CACHR_BUCKET")
var AWS_REGION = os.Getenv("AWS_REGION")

func main() {

	if CACHR_BUCKET == "" {
		log.Fatalf("missing required env: CACHR_BUCKET")
	}

	if AWS_REGION == "" {
		log.Fatalf("missing required env: AWS_REGION")
	}

	if len(os.Args) < 2 {
		log.Fatalf("Usage: cachr <commandname> <args>")
	}

	command := os.Args[1]
	args := os.Args[2:]

	session, err := session.NewSession(&aws.Config{
		Region: aws.String(AWS_REGION)},
	)
	if err != nil {
		log.Fatalf("cannot create new AWS session: %v", err)
	}

	if command == "exists" {
		if len(args) < 1 {
			log.Fatalf("Usage: cachr exists <keyname>")
		}

		exists, err := keyExists(session, CACHR_BUCKET, args[0])
		if err != nil {
			log.Fatalf("failed to check if key exists: %v", err)
		}

		if exists {
			os.Exit(0)
		}

		os.Exit(1)

	} else if command == "get" {
		if len(args) < 1 {
			log.Fatalf("Usage: cachr get <keyname>")
		}

		err = getKey(session, CACHR_BUCKET, args[0])
		if err != nil {
			log.Fatalf("failed to get cache: %v", err)
		}

	} else if command == "save" {
		if len(args) < 2 {
			log.Fatalf("Usage: cachr save <keyname> <files>...")
		}

		err = saveKey(session, CACHR_BUCKET, args[0], args[1:])
		if err != nil {
			log.Fatalf("failed to save cache: %v", err)
		}

	} else {
		log.Fatalf("invalid command: %s", command)
	}
}
