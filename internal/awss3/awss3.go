package awss3

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"

	"git.cto.ai/provision/internal/logger"
	"git.cto.ai/sdk-go/pkg/sdk"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func EBS3Setup(sdk *sdk.Sdk, awsSess *session.Session, unzippedRepo, awsRegion string) (string, error) {
	s3Client := s3.New(awsSess, aws.NewConfig().WithRegion(awsRegion))
	s3UploaderClient := s3manager.NewUploader(awsSess)

	bucketName, err := createBucket(sdk, s3Client, unzippedRepo, 0)
	if err != nil {
		return bucketName, err
	}

	err = uploadZip(sdk, s3UploaderClient, awsRegion, bucketName, unzippedRepo)
	if err != nil {
		return bucketName, err
	}

	return bucketName, nil
}

func createBucket(sdk *sdk.Sdk, svc *s3.S3, unzippedRepoName string, retries int) (string, error) {
	time := time.Now()
	bucketName := fmt.Sprintf("%s-%v", strings.ToLower(unzippedRepoName), time.Format("20060102150405"))
	logger.LogSlack(sdk, "🔄 Creating S3 bucket...")

	input := &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	}

	_, err := svc.CreateBucket(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return "", aerr
		}
		return "", err
	}

	logger.LogSlack(sdk, "✅ S3 bucket created.")
	return bucketName, err
}

func uploadZip(sdk *sdk.Sdk, svc *s3manager.Uploader, awsRegion, bucketName, targetFile string) error {
	logger.LogSlack(sdk, "🔄 Uploading repository files to S3 bucket...")

	filename := fmt.Sprintf("%s.zip", targetFile)

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = svc.Upload(&s3manager.UploadInput{
		Bucket: aws.String(strings.ToLower(bucketName)),
		Key:    aws.String(filepath.Base(filename)),
		Body:   file,
	})
	if err != nil {
		return err
	}

	logger.LogSlack(sdk, fmt.Sprintf("✅ S3 Bucket: https://s3.console.aws.amazon.com/s3/buckets/%s/?region=%s&tab=overview", bucketName, awsRegion))
	return nil
}
