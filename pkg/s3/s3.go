package s3

import (
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/internal/infra/env"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/log"
)

type S3Interface interface {
	Upload(file *multipart.FileHeader) (string, error)
}

type S3Struct struct {
	session  *session.Session
	uploader *s3manager.Uploader
}

var S3 = getS3()

func getS3() S3Interface {
	session := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(env.AppEnv.AWSRegion),
		Credentials: credentials.NewStaticCredentials(
			env.AppEnv.AWSAccessKeyID,
			env.AppEnv.AWSSecretAccessKey,
			"",
		),
	}))

	uploader := s3manager.NewUploader(session)

	return &S3Struct{
		session:  session,
		uploader: uploader,
	}
}

func (s *S3Struct) Upload(file *multipart.FileHeader) (string, error) {
	// Validate input
	if file == nil {
		err := errors.New("file is nil")
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[S3][Upload] invalid input")
		return "", err
	}

	// Open file
	fileContent, err := file.Open()
	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[S3][Upload] failed to open file")
		return "", err
	}
	defer func() {
		if cerr := fileContent.Close(); cerr != nil {
			log.Warn(log.LogInfo{
				"error": cerr.Error(),
			}, "[S3][Upload] failed to close file")
		}
	}()

	// Generate unique file name
	timeNow := time.Now().Unix()
	fileName := fmt.Sprintf("%d-%s", timeNow, file.Filename)

	// Determine content type with fallback
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream" // Fallback for unknown file types
	}

	// Upload to S3
	result, err := s.uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(env.AppEnv.AWSS3BucketName),
		Key:         aws.String(fileName),
		Body:        fileContent,
		ACL:         aws.String("public-read"),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[S3][Upload] failed to upload file")
		return "", err
	}

	// Return public URL
	return result.Location, nil
}
