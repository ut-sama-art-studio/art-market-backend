package fileservice

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"

	"github.com/99designs/gqlgen/graphql"
	"github.com/disintegration/imaging"
)

// CheckAndResizeImage checks if the image is larger than 1MB, and resizes it if necessary.
func CheckAndResizeImage(file graphql.Upload) (io.Reader, error) {
	// Create a buffer to hold file content
	var buf bytes.Buffer
	size, err := io.Copy(&buf, file.File)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// If the file size is within the limit, return as is
	if size <= int64(maxFileSize) {
		return &buf, nil
	}

	// Decode and resize if it exceeds max size
	img, _, err := image.Decode(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Resize the image to fit within the max size
	resizedImg := imaging.Resize(img, 3000, 0, imaging.Lanczos) // Width 800 and proportional height
	var resizedBuf bytes.Buffer
	err = jpeg.Encode(&resizedBuf, resizedImg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encode resized image: %w", err)
	}

	return &resizedBuf, nil
}
