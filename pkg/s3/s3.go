package s3

import (
	"mime/multipart"

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
	log.Info(log.LogInfo{
		"message": "Creating new S3 session",
		"region":  env.AppEnv.AWSRegion,
		"bucket":  env.AppEnv.AWSS3BucketName,
		"access":  env.AppEnv.AWSAccessKeyID,
		"secret":  env.AppEnv.AWSSecretAccessKey,
	}, "[S3][getS3] creating new S3 session")

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
	fileContent, err := file.Open()
	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[S3][Upload] failed to open file")

		return "", err
	}

	defer fileContent.Close()

	fileName := file.Filename

	val, err := s.session.Config.Credentials.Get()
	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[S3][Upload] failed to get credentials")
	}

	log.Info(log.LogInfo{
		"access": val.AccessKeyID,
		"secret": val.SecretAccessKey,
		"token":  val.SessionToken,
		"provider":   val.ProviderName,
	}, "[S3][Upload] uploading file")

	result, err := s.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(env.AppEnv.AWSS3BucketName),
		Key:    aws.String(fileName),
		Body:   fileContent,
	})
	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[S3][Upload] failed to upload file")
	}

	return result.Location, nil
}
