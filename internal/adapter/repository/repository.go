package repository

import (
	"gorm.io/gorm"

	"parser/internal/model"
)

type DomainRepository struct {
	db *gorm.DB
}

func NewDomainRepository(db *gorm.DB) *DomainRepository {
	return &DomainRepository{db: db}
}

func (r *DomainRepository) Create(domain *model.Domain) error {
	return r.db.Create(domain).Error
}

func (r *DomainRepository) Save(domain *model.Domain) error {
	return r.db.Save(domain).Error
}

func (r *DomainRepository) GetAll() ([]model.Domain, error) {
	var domains []model.Domain
	err := r.db.Order("created_at DESC").Find(&domains).Error
	return domains, err
}