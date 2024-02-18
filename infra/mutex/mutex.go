package mutex

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/guilhermealvess/guicpay/domain/gateway"
)

type mutex struct {
	client *redis.Client
}

func NewMutex(address, password string) gateway.Mutex {
	return &mutex{
		client: redis.NewClient(&redis.Options{
			Addr:     address,
			Password: password,
			DB:       0,
		}),
	}
}

func (m *mutex) Lock(ctx context.Context, key string, ttl time.Duration) error {
	ok, err := m.client.SetNX("MUTEX::"+key, "LOCK", ttl).Result()
	if !ok {
		return fmt.Errorf("mutex: key is locked. %w", err)
	}

	return nil
}

func (m *mutex) Unlock(ctx context.Context, key string) error {
	_, err := m.client.Del("MUTEX::" + key).Result()
	if err != nil {
		return fmt.Errorf("mutex: %w", err)
	}

	return nil
}
