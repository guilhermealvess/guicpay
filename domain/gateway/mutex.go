package gateway

import (
	"context"
	"time"
)

type Mutex interface {
	Lock(ctx context.Context, key string, ttl time.Duration) error
	Unlock(ctx context.Context, key string) error
}
