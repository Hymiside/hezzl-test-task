package service

import (
	"context"
	"fmt"
	"github.com/Hymiside/hezzl-test-task/pkg/natsqueue"
	"github.com/Hymiside/hezzl-test-task/pkg/repository/postgres"
	"github.com/Hymiside/hezzl-test-task/pkg/repository/redis"

	"github.com/Hymiside/hezzl-test-task/pkg/models"
)

type Service struct {
	postgres *postgres.Repository
	redis    *redis.Repository
	nc       *natsqueue.Nats
}

func NewService(r *postgres.Repository, redis *redis.Repository, nc *natsqueue.Nats) *Service {
	go nc.Sub()
	return &Service{postgres: r, redis: redis, nc: nc}
}

func (s *Service) CreateItem(ctx context.Context, ni models.NewItem) (models.Item, error) {
	itemId, priority, err := s.postgres.CreateItem(ctx, ni)
	if err != nil {
		return models.Item{}, fmt.Errorf("failed to create item in postgres: %w", err)
	}

	item := models.Item{
		ID:          itemId,
		CampaignId:  ni.CampaignId,
		Name:        ni.Name,
		Description: ni.Description,
		Priority:    priority,
		Removed:     ni.Removed,
		CreatedAt:   ni.CreatedAt,
	}

	if err = s.UpdateItemsRedis(ctx); err != nil {
		return models.Item{}, err
	}

	s.nc.Pub([]byte(fmt.Sprintf("SUCCESSFUL-create new item-%d", itemId)))
	return item, nil
}

func (s *Service) GetAllItems(ctx context.Context) ([]models.Item, error) {
	items, err := s.redis.GetAllItems(ctx)
	if err == nil {
		return items, nil
	}

	items, err = s.postgres.GetAllItems(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get items from postgres: %w", err)
	}

	if err = s.redis.SetItem(ctx, items); err != nil {
		return nil, fmt.Errorf("failed to add items in redis: %w", err)
	}
	return items, nil
}

func (s *Service) UpdateItem(ctx context.Context, campaignId, itemId int, name, description string) (models.Item, error) {
	item, err := s.postgres.UpdateItem(ctx, campaignId, itemId, name, description)
	if err != nil {
		return models.Item{}, fmt.Errorf("failed to update item in postgres: %w", err)
	}

	if err = s.UpdateItemsRedis(ctx); err != nil {
		return models.Item{}, err
	}

	s.nc.Pub([]byte(fmt.Sprintf("SUCCESSFUL-update item-%d", itemId)))
	return item, nil
}

func (s *Service) DeleteItem(ctx context.Context, campaignId, itemId int) error {
	if err := s.postgres.DeleteItem(ctx, campaignId, itemId); err != nil {
		return fmt.Errorf("failed to delete item in postgres: %w", err)
	}
	if err := s.redis.DeleteItem(ctx); err != nil {
		return fmt.Errorf("failed to delete item in redis: %w", err)
	}
	s.nc.Pub([]byte(fmt.Sprintf("SUCCESSFUL-delete item-%d", itemId)))
	return nil
}

func (s *Service) UpdateItemsRedis(ctx context.Context) error {
	if err := s.redis.DeleteItem(ctx); err != nil {
		return fmt.Errorf("failed to update item in redis: %w", err)
	}

	items, err := s.postgres.GetAllItems(ctx)
	if err != nil {
		return fmt.Errorf("failed to get items from postgres: %w", err)
	}

	if err = s.redis.SetItem(ctx, items); err != nil {
		return fmt.Errorf("failed to update item in redis: %w", err)
	}
	return nil
}
