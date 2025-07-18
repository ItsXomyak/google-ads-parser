package service

import (
	"parser/internal/adapter/repository"
	"parser/internal/model"
)

type DomainService struct {
	repo *repository.DomainRepository
}

func NewDomainService(repo *repository.DomainRepository) *DomainService {
	return &DomainService{repo: repo}
}

func (s *DomainService) Save(domain *model.Domain) error {
	return s.repo.Save(domain)
}

//GetAll
func (s *DomainService) GetAll() ([]model.Domain, error) {
	return s.repo.GetAll()
}