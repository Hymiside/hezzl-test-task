package rediscache

import (
	"context"
	"fmt"
	"time"

	"github.com/Hymiside/hezzl-test-task/pkg/models"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

type Redis struct {
	rdb *redis.Client
	ch  *cache.Cache
}

type ConfigRedis struct {
	Host string
	Port string
}

func NewRedis(c ConfigRedis) (*Redis, error) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", c.Host, c.Port),
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	ch := cache.New(&cache.Options{
		Redis: rdb,
	})
	return &Redis{rdb: rdb, ch: ch}, nil
}

func (c *Redis) CloseRedis() error {
	if err := c.rdb.Close(); err != nil {
		return err
	}
	return nil
}

func (c *Redis) GetItems() ([]models.Item, error) {
	var i []models.Item
	ctx := context.Background()

	if err := c.ch.Get(ctx, "DataItems", i); err != nil {
		return nil, err
	}
	return i, nil
}

func (c *Redis) SetItems(i []models.Item) error {
	ctx := context.Background()

	if err := c.ch.Set(&cache.Item{
		Ctx:   ctx,
		Key:   "DataItems",
		Value: i,
		TTL:   time.Minute,
	}); err != nil {
		return err
	}
	return nil
}
