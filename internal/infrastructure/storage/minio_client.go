package storage

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	appconfig "github.com/hajimohammadinet/dabir/internal/config"
	domainstorage "github.com/hajimohammadinet/dabir/internal/domain/storage"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOClient struct {
	client *minio.Client
	bucket string
}

func NewMinIOClient(cfg appconfig.StorageConfig) (*MinIOClient, error) {
	endpoint := normalizeEndpoint(cfg.Endpoint)

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
		Region: cfg.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	return &MinIOClient{
		client: client,
		bucket: cfg.Bucket,
	}, nil
}

func (c *MinIOClient) PutObject(ctx context.Context, input domainstorage.PutObjectInput) error {
	_, err := c.client.PutObject(
		ctx,
		c.bucket,
		input.ObjectKey,
		input.Reader,
		input.SizeBytes,
		minio.PutObjectOptions{
			ContentType: input.ContentType,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to put object: %w", err)
	}

	return nil
}

func (c *MinIOClient) PresignedGetObjectURL(ctx context.Context, input domainstorage.PresignedGetObjectInput) (string, error) {
	ttl := input.TTL
	if ttl <= 0 {
		ttl = 15 * time.Minute
	}

	presignedURL, err := c.client.PresignedGetObject(
		ctx,
		c.bucket,
		input.ObjectKey,
		ttl,
		url.Values{},
	)
	if err != nil {
		return "", fmt.Errorf("failed to create presigned get object url: %w", err)
	}

	return presignedURL.String(), nil
}

func (c *MinIOClient) DeleteObject(ctx context.Context, objectKey string) error {
	if strings.TrimSpace(objectKey) == "" {
		return nil
	}

	if err := c.client.RemoveObject(
		ctx,
		c.bucket,
		objectKey,
		minio.RemoveObjectOptions{},
	); err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	return nil
}

func normalizeEndpoint(endpoint string) string {
	endpoint = strings.TrimSpace(endpoint)
	endpoint = strings.TrimPrefix(endpoint, "http://")
	endpoint = strings.TrimPrefix(endpoint, "https://")

	return endpoint
}
