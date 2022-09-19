package service

import (
	"github.com/Hymiside/hezzl-test-task/pkg/models"
	"github.com/Hymiside/hezzl-test-task/pkg/rediscache"
	"github.com/Hymiside/hezzl-test-task/pkg/repository"
	"log"
)

type Service struct {
	repo *repository.Repository
	ch   *rediscache.Redis
}

func NewService(r repository.Repository, ch rediscache.Redis) *Service {
	return &Service{repo: &r, ch: &ch}
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
	dataItems, err := r.GetItemsRedis()
	if err != nil {
		dataItems, err = r.repo.GetItems()
		if err != nil {
			return nil, err
		}

		if err = r.SetItemsRedis(dataItems); err != nil {
			log.Println("error set cache") // необходимо отправить лог в ClickHouse
		}
		return dataItems, nil
	}
	return dataItems, nil
}

func (r *Service) GetItemsRedis() ([]models.Item, error) {
	dataCache, err := r.ch.GetItems()
	if err != nil {
		return nil, err
	}
	return dataCache, err
}

func (r *Service) SetItemsRedis(i []models.Item) error {
	if err := r.ch.SetItems(i); err != nil {
		return err
	}
	return nil
}
