package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/cache/v8"

	"github.com/go-redis/redis/v8"

	"github.com/Hymiside/hezzl-test-task/pkg/models"
)

type ConfigRedis struct {
	Host string
	Port string
}

type Repository struct {
	redis *redis.Client
	ch    *cache.Cache
}

const ttl = time.Minute

func NewRepository(ctx context.Context, c ConfigRedis) (*Repository, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", c.Host, c.Port),
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}
	ch := cache.New(&cache.Options{
		Redis: rdb,
	})

	go func(ctx context.Context) {
		<-ctx.Done()
		rdb.Close()
	}(ctx)

	return &Repository{redis: rdb, ch: ch}, nil
}

func (r *Repository) GetAllItems(ctx context.Context) ([]models.Item, error) {
	var items []models.Item
	if err := r.ch.Get(ctx, "DataItems", &items); err != nil {
		return nil, err
	}
	if items == nil {
		return nil, errors.New("not found items in redis")
	}
	return items, nil
}

func (r *Repository) SetItem(ctx context.Context, items []models.Item) error {
	if err := r.ch.Set(&cache.Item{
		Ctx:   ctx,
		Key:   "DataItems",
		Value: items,
		TTL:   ttl,
	}); err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteItem(ctx context.Context) error {
	if err := r.ch.Delete(ctx, "DataItems"); err != nil {
		return err
	}
	return nil
}
