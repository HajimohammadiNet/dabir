package settings

import "context"

type Repository interface {
	Set(ctx context.Context, key string, value []byte) error
	Get(ctx context.Context, key string) (*Setting, error)
}
