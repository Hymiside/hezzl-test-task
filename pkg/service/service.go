package service

import (
	"log"

	"github.com/Hymiside/hezzl-test-task/pkg/natsqueue"

	"github.com/Hymiside/hezzl-test-task/pkg/models"
	"github.com/Hymiside/hezzl-test-task/pkg/rediscache"
	"github.com/Hymiside/hezzl-test-task/pkg/repository"
)

type Service struct {
	repo *repository.Repository
	ch   *rediscache.Redis
	nc   *natsqueue.Nats
}

func NewService(r repository.Repository, ch rediscache.Redis, nc natsqueue.Nats) *Service {
	return &Service{repo: &r, ch: &ch, nc: &nc}
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

func (r *Service) DeleteItemsRedis() error {
	if err := r.ch.DeleteItems(); err != nil {
		return err
	}
	return nil
}

func (r *Service) UpdateItem(i models.Item) ([]models.Item, error) {
	var dataItems []models.Item

	ui, err := r.repo.UpdateItem(i)
	if err != nil {
		return nil, err
	}

	if err = r.DeleteItemsRedis(); err != nil {
		return nil, err
	}
	dataItems, err = r.repo.GetItems()
	if err != nil {
		return nil, err
	}
	if err = r.SetItemsRedis(dataItems); err != nil {
		log.Println("error set cache") // необходимо отправить лог в ClickHouse
	}
	return ui, nil
}

func (r *Service) DeleteItem(i models.Item) ([]models.Item, error) {
	var dataItems []models.Item

	deleteItem, err := r.repo.DeleteItem(i)
	if err != nil {
		return nil, err
	}

	if err = r.DeleteItemsRedis(); err != nil {
		return nil, err
	}
	dataItems, err = r.repo.GetItems()
	if err != nil {
		return nil, err
	}
	if err = r.SetItemsRedis(dataItems); err != nil {
		log.Println("error set cache") // необходимо отправить лог в ClickHouse
	}
	return deleteItem, nil
}
