package fileservice

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/99designs/gqlgen/graphql"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
)

var (
	bucketName = "art-market-bucket"
	endpoint   = "tor1.digitaloceanspaces.com"
	region     = "tor1"
	s3Client   *s3.S3
	uploader   *s3manager.Uploader
)

// Initialize S3 Client
func InitS3Client() error {
	accessKey := os.Getenv("DO_BUCKET_ACCESS_KEY")
	secretKey := os.Getenv("DO_BUCKET_SECRET_KEY")

	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:    &endpoint,
		Region:      &region,
	}))

	s3Client = s3.New(sess)

	uploader = s3manager.NewUploader(sess)
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
		filePath = fmt.Sprintf("%s/%s/%s", bucketName, userId, file.Filename)
	} else {
		filePath = fmt.Sprintf("%s/%s/%s/%s", bucketName, userId, folderPath,
			file.Filename)
	}

	// Read the file content into a buffer
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, resizedFile); err != nil {
		return "", fmt.Errorf("failed to read file content: %w", err)
	}

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(filePath),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(file.ContentType),
		ACL:         aws.String("public-read"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return result.Location, nil
}

// DeleteUserFolder deletes all objects in the specified user's folder
func DeleteUserFolder(userId string, folderPath string) error {
	if folderPath != "" {
		folderPath = fmt.Sprintf("%s/%s/%s/", bucketName, userId, folderPath)
	}
	// List all objects in the folder
	listOutput, err := s3Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(folderPath),
	})
	if err != nil {
		fmt.Printf("Error listing objects: %s\n", err)
		return err
	}
	// Delete each object
	for _, obj := range listOutput.Contents {
		_, err = s3Client.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(bucketName),
			Key:    obj.Key,
		})
		if err != nil {
			return fmt.Errorf("failed to delete object in folder %s: %w", folderPath, err)
		}
	}
	return nil
}

// LogBucketObjects lists and logs all objects in the S3 bucket
func LogBucketObjects() error {
	output, err := s3Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		fmt.Printf("Error listing objects: %s\n", err)
		return err
	}
	fmt.Println("Objects in bucket:")
	for _, object := range output.Contents {
		fmt.Printf("Name=%s, Size=%d\n", aws.StringValue(object.Key), aws.Int64Value(object.Size))
	}
	return nil
}
