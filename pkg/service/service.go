package service

import (
	"github.com/Hymiside/hezzl-test-task/pkg/models"
	"github.com/Hymiside/hezzl-test-task/pkg/repository"
)

type Service struct {
	repo *repository.Repository
}

func NewService(r repository.Repository) *Service {
	return &Service{repo: &r}
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
