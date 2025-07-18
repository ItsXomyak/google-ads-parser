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

func (r *DomainRepository) DeleteByID(id uint) error {
	return r.db.Delete(&model.Domain{}, id).Error
}

func (r *DomainRepository) GetPaginated(domain string, limit, offset int) ([]model.Domain, int64, error) {
	var domains []model.Domain
	var total int64

	query := r.db.Model(&model.Domain{})

	if domain != "" {
		query = query.Where("domain ILIKE ?", "%"+domain+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&domains).Error

	return domains, total, err
}
