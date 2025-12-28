package storage

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// R2Client represents a client for interacting with Cloudflare R2.
type R2Client struct {
	S3Client   *s3.Client
	BucketName string
}

// NewR2Client creates a new R2 client.
func NewR2Client() (*R2Client, error) {
	accountID := os.Getenv("R2_ACCOUNT_ID")
	accessKeyID := os.Getenv("R2_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("R2_SECRET_ACCESS_KEY")
	bucketName := os.Getenv("R2_BUCKET_NAME")

	if accountID == "" || accessKeyID == "" || accessKeySecret == "" || bucketName == "" {
		return nil, fmt.Errorf("R2 environment variables are not fully set")
	}

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, accessKeySecret, "")),
	)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg)

	return &R2Client{
		S3Client:   s3Client,
		BucketName: bucketName,
	}, nil
}

// UploadChunk uploads a chunk to R2.
func (c *R2Client) UploadChunk(ctx context.Context, chunkHash string, data io.Reader) error {
	_, err := c.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &c.BucketName,
		Key:    &chunkHash,
		Body:   data,
	})
	return err
}

// DownloadChunk downloads a chunk from R2.
func (c *R2Client) DownloadChunk(ctx context.Context, chunkHash string) (io.ReadCloser, error) {
	output, err := c.S3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &c.BucketName,
		Key:    &chunkHash,
	})
	if err != nil {
		return nil, err
	}
	return output.Body, nil
}
