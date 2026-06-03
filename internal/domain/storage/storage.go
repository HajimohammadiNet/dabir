package storage

import (
	"context"
	"io"
	"time"
)

type PutObjectInput struct {
	ObjectKey   string
	Reader      io.Reader
	SizeBytes   int64
	ContentType string
}

type PresignedGetObjectInput struct {
	ObjectKey string
	TTL       time.Duration
}

type Client interface {
	PutObject(ctx context.Context, input PutObjectInput) error
	PresignedGetObjectURL(ctx context.Context, input PresignedGetObjectInput) (string, error)
	DeleteObject(ctx context.Context, objectKey string) error
}
