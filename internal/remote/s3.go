package remote

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/ms-choudhary/cachr"
)

func Exists(bucket, key string) (bool, error) {
	session := cachr.AWSSession()
	svc := s3.New(session)

	results, err := svc.ListObjects(&s3.ListObjectsInput{Bucket: aws.String(bucket), Prefix: aws.String(key)})
	if err != nil {
		return false, err
	}

	return len(results.Contents) > 0, nil
}

func Upload(filename, bucket, key string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("cannot open file for saving: %v", err)
	}

	defer file.Close()

	session := cachr.AWSSession()
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

func Download(bucket, key) error {
	session := cachr.AWSSession()
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
