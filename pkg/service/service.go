package service

import "github.com/Hymiside/hezzl-test-task/pkg/repository"

type Service struct {
	repo *repository.Repository
}

func NewService(r repository.Repository) *Service {
	return &Service{repo: &r}
}
