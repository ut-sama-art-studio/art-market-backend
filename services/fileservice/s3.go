package fileservice

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// Global S3 client and configuration
var (
	maxFileSize = 4 * 1024 * 1024 // 4MB in bytes

	bucketName = "utsama-art-market"
	region     = "us-east-1"
	s3Client   *s3.Client
	uploader   *manager.Uploader
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

// UploadFileToS3 handles the file upload to S3 and returns the file's URL
func UploadFileToS3(file graphql.Upload, userId string, folderPath string) (string, error) {
	// Limit file size, resize if over
	resizedFile, err := CheckAndResizeImage(file)
	if err != nil {
		return "", fmt.Errorf("file processing error: %w", err)
	}

	// Use UUID as filename
	file.Filename = uuid.New().String()

	// Set the file key to include the user ID, used to access this file later
	var filePath string
	if folderPath == "" {
		filePath = fmt.Sprintf("%s/%s", userId, file.Filename)
	} else {
		filePath = fmt.Sprintf("%s/%s/%s", userId, folderPath, file.Filename)
	}

	input := &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(filePath),
		Body:        resizedFile,
		ContentType: aws.String(file.ContentType),
	}

	_, err = uploader.Upload(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Construct file URL
	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, region, filePath)

	return fileURL, nil
}

func EmptyUserFolder(userId string, folderPath string) error {
	if folderPath != "" {
		folderPath = fmt.Sprintf("%s/%s", userId, folderPath)
	}

	// List all objects in the specified folder
	output, err := s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(folderPath + "/"), // Ensure the prefix ends with "/"
	})
	if err != nil {
		return fmt.Errorf("failed to list objects: %w", err)
	}

	if len(output.Contents) == 0 {
		return nil
	}

	// Delete each object found in the folder
	for _, object := range output.Contents {
		_, err := s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
			Bucket: aws.String(bucketName),
			Key:    object.Key,
		})
		if err != nil {
			return fmt.Errorf("failed to delete object %s: %w", *object.Key, err)
		}
	}
	return nil
}

func LogBucketObjects() error {
	output, err := s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
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
