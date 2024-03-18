package utils

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

func UploadImg(file *multipart.FileHeader, c *gin.Context) string {

	bucketName := "golang_twitter"

	client, err := createGCSClient(c)
	if err != nil {
		log.Fatal(err)
	}

	currentTime := time.Now()
	gcsFileName := fmt.Sprintf("%s.png", currentTime.Format("20060102150405"))

	src, err := file.Open()
	if err != nil {
		return "error is occurred while file opening"
	}
	defer src.Close()

	bucket := client.Bucket(bucketName)
	obj := bucket.Object(gcsFileName)

	wc := obj.NewWriter(c)
	if _, err = io.Copy(wc, src); err != nil {
		return "error is occurred while file copying"
	}
	if err = wc.Close(); err != nil {
		return "error is occurred while file closing"
	}

	resImagePath := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, gcsFileName)

	return resImagePath
}

func createGCSClient(ctx *gin.Context) (*storage.Client, error) {
	credentialFilePath := "./gcp.json"
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialFilePath))
	if err != nil {
		return nil, fmt.Errorf("failed to create client : %w", err)
	}
	defer client.Close()
	return client, err
}
