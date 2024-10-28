package files

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Global S3 client and configuration
var (
	s3BucketName = "utsama-art-market"
	region       = "us-east-1"
	s3Client     *s3.Client
	uploader     *manager.Uploader
)

// Initialize S3 Client
func InitS3Client() error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}
	s3Client = s3.NewFromConfig(cfg)
	uploader = manager.NewUploader(s3Client)
	return nil
}

func LogBucketObjects() error {
	output, err := s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(s3BucketName),
	})
	if err != nil {
		fmt.Printf("Error listing objects: %s\n", err)
		return err
	}
	fmt.Println("Objects in bucket:")
	for _, object := range output.Contents {
		fmt.Printf("Name=%s, Size=%d\n", aws.ToString(object.Key), object.Size)
	}
	return nil
}

// UploadFileToS3 handles the file upload to S3 and returns the file's URL
func UploadFileToS3(file graphql.Upload, userId string) (string, error) {
	// Get file type from the header
	contentType := file.ContentType

	// Set the file key to include the user ID, used to access this file later
	fileKey := fmt.Sprintf("%s/%s", userId, file.Filename)

	input := &s3.PutObjectInput{
		Bucket:      aws.String(s3BucketName),
		Key:         aws.String(fileKey),
		Body:        file.File,
		ContentType: aws.String(contentType),
	}

	_, err := uploader.Upload(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Construct file URL
	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s3BucketName, region, file.Filename)

	return fileURL, nil
}
