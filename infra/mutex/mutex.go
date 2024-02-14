package mutex

import (
	"context"
	"errors"
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
		return errors.New("mutex: key is locked")
	}

	return fmt.Errorf("mutex: %w", err)
}

func (m *mutex) Unlock(ctx context.Context, key string) error {
	_, err := m.client.Del("MUTEX::" + key).Result()
	return fmt.Errorf("mutex: %w", err)
}
