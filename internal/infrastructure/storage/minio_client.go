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
	client        *minio.Client
	presignClient *minio.Client
	bucket        string
}

func NewMinIOClient(cfg appconfig.StorageConfig) (*MinIOClient, error) {
	client, err := newClient(cfg.Endpoint, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	presignClient := client
	if strings.TrimSpace(cfg.PublicEndpoint) != "" {
		presignClient, err = newClient(cfg.PublicEndpoint, cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create public minio client: %w", err)
		}
	}

	return &MinIOClient{
		client:        client,
		presignClient: presignClient,
		bucket:        cfg.Bucket,
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

	presignedURL, err := c.presignClient.PresignedGetObject(
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

func newClient(endpoint string, cfg appconfig.StorageConfig) (*minio.Client, error) {
	secure := cfg.UseSSL
	trimmedEndpoint := strings.TrimSpace(endpoint)

	if parsed, err := url.Parse(trimmedEndpoint); err == nil && parsed.Scheme != "" {
		secure = strings.EqualFold(parsed.Scheme, "https")
	}

	bucketLookup := minio.BucketLookupAuto
	if cfg.ForcePathStyle {
		bucketLookup = minio.BucketLookupPath
	}

	return minio.New(normalizeEndpoint(trimmedEndpoint), &minio.Options{
		Creds:        credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure:       secure,
		Region:       cfg.Region,
		BucketLookup: bucketLookup,
	})
}
