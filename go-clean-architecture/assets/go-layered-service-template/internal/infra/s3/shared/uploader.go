package shared

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

// Uploader is an AWS S3 uploader implementation.
type Uploader struct {
	bucket string
	client *awss3.Client
}

// NewUploader creates a new S3 uploader using the AWS SDK v2 default config chain.
func NewUploader(ctx context.Context, bucket string) (*Uploader, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	return &Uploader{
		bucket: bucket,
		client: awss3.NewFromConfig(cfg),
	}, nil
}

// PutObject uploads an object to S3.
func (u *Uploader) PutObject(ctx context.Context, key string, body io.Reader) error {
	_, err := u.client.PutObject(ctx, &awss3.PutObjectInput{
		Bucket: aws.String(u.bucket),
		Key:    aws.String(key),
		Body:   body,
	})
	return err
}
