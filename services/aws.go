package services

import (
	"log"

	"github.com/EgSundqvist/imagebook-imageapi/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// InitAWSSession initializes an AWS session
func InitAWSSession() (*session.Session, error) {
	config.LoadConfig()

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.AppConfig.AWSIAMUser.Region),
		Credentials: credentials.NewStaticCredentials(
			config.AppConfig.AWSIAMUser.AccessKeyID,
			config.AppConfig.AWSIAMUser.SecretAccessKey,
			"",
		),
	})
	if err != nil {
		log.Fatalf("Failed to create AWS session: %v", err)
		return nil, err
	}

	return sess, nil
}

// GetS3Object retrieves an object from S3
func GetS3Object(sess *session.Session, bucket, key string) (*s3.GetObjectOutput, error) {
	svc := s3.New(sess)
	resp, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
