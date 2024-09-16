package aws

import (
	"context"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Service interface {
	GetAllFiles()
	UploadFile(string, io.Reader) (string, error)
}

type service struct {
	s3 *s3.Client
}

func New() Service {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-2"))
	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)

	return &service{
		s3: client,
	}
}

func (s *service) UploadFile(fileName string, file io.Reader) (string, error) {
	uploader := manager.NewUploader(s.s3)
	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("buntservers3bucket"),
		Key:    aws.String(fileName),
		Body:   file,
	})

	if err != nil {
		return "", err
	} else {
		return result.Location, err
	}
}

func (s *service) GetAllFiles() {
	// Get the first page of results for ListObjectsV2 for a bucket
	output, err := s.s3.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String("buntservers3bucket"),
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("first page results:")
	for _, object := range output.Contents {
		log.Printf("key=%s size=%d", aws.ToString(object.Key), object.Size)
	}
}
