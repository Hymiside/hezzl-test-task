package service

import (
	"github.com/Hymiside/hezzl-test-task/pkg/models"
	"github.com/Hymiside/hezzl-test-task/pkg/rediscache"
	"github.com/Hymiside/hezzl-test-task/pkg/repository"
	"log"
)

type Service struct {
	repo *repository.Repository
}

type RedisCache struct {
	ch *rediscache.Redis
}

func NewService(r repository.Repository, ch rediscache.Redis) (*Service, *RedisCache) {
	return &Service{repo: &r}, &RedisCache{ch: &ch}
}

func (r *Service) CreateItem(ni models.NewItem) (models.Item, error) {
	if err := r.repo.GetCampaign(ni); err != nil {
		return models.Item{}, err
	}

	i, err := r.repo.CreateItem(ni)
	if err != nil {
		return models.Item{}, err
	}
	return i, nil
}

func (r *Service) GetItems() ([]models.Item, error) {
	var c *RedisCache

	dataItems, err := c.GetItemsRedis()
	if err != nil {

		dataItems, err = r.repo.GetItems()
		if err != nil {
			return nil, err
		}

		if err = c.SetItemsRedis(dataItems); err != nil {
			log.Println("error set cache") // необходимо отправить лог в ClickHouse
		}
		return dataItems, nil
	}
	return dataItems, nil
}

func (c *RedisCache) GetItemsRedis() ([]models.Item, error) {
	dataCache, err := c.ch.GetItems()
	if err != nil {
		return nil, err
	}
	return dataCache, err
}

func (c *RedisCache) SetItemsRedis(i []models.Item) error {
	if err := c.ch.SetItems(i); err != nil {
		return err
	}
	return nil
}
