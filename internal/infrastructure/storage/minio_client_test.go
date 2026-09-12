package storage

import (
	"context"
	"net/url"
	"testing"
	"time"

	appconfig "github.com/hajimohammadinet/dabir/internal/config"
	domainstorage "github.com/hajimohammadinet/dabir/internal/domain/storage"
)

func TestPresignedURLUsesPublicEndpoint(t *testing.T) {
	client, err := NewMinIOClient(appconfig.StorageConfig{
		Endpoint:       "http://minio:9000",
		PublicEndpoint: "http://localhost:9000",
		Region:         "us-east-1",
		Bucket:         "dabir-attachments",
		AccessKey:      "dabir",
		SecretKey:      "dabir_minio_secret",
		ForcePathStyle: true,
	})
	if err != nil {
		t.Fatalf("create client: %v", err)
	}

	presigned, err := client.PresignedGetObjectURL(context.Background(), domainstorage.PresignedGetObjectInput{
		ObjectKey: "letters/letter-id/scan.pdf",
		TTL:       15 * time.Minute,
	})
	if err != nil {
		t.Fatalf("presign object: %v", err)
	}

	parsed, err := url.Parse(presigned)
	if err != nil {
		t.Fatalf("parse presigned url: %v", err)
	}
	if parsed.Scheme != "http" || parsed.Host != "localhost:9000" {
		t.Fatalf("presigned endpoint = %s://%s", parsed.Scheme, parsed.Host)
	}
	if parsed.Path != "/dabir-attachments/letters/letter-id/scan.pdf" {
		t.Fatalf("presigned path = %q", parsed.Path)
	}
	if parsed.Query().Get("X-Amz-Signature") == "" {
		t.Fatal("presigned url is missing its signature")
	}
}

func TestPresignedURLFallsBackToInternalEndpoint(t *testing.T) {
	client, err := NewMinIOClient(appconfig.StorageConfig{
		Endpoint:       "minio:9000",
		Region:         "us-east-1",
		Bucket:         "dabir-attachments",
		AccessKey:      "dabir",
		SecretKey:      "dabir_minio_secret",
		ForcePathStyle: true,
	})
	if err != nil {
		t.Fatalf("create client: %v", err)
	}

	presigned, err := client.PresignedGetObjectURL(context.Background(), domainstorage.PresignedGetObjectInput{
		ObjectKey: "scan.pdf",
		TTL:       time.Minute,
	})
	if err != nil {
		t.Fatalf("presign object: %v", err)
	}

	parsed, err := url.Parse(presigned)
	if err != nil {
		t.Fatalf("parse presigned url: %v", err)
	}
	if parsed.Host != "minio:9000" {
		t.Fatalf("presigned host = %q", parsed.Host)
	}
}
