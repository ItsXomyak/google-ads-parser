package repository

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"parser/internal/model"
)

type DomainRepository struct {
	db *gorm.DB
}

func NewDomainRepository(db *gorm.DB) *DomainRepository {
	return &DomainRepository{db: db}
}

// isErrorData проверяет, содержат ли данные ошибку парсинга
func (r *DomainRepository) isErrorData(legalName, country string) bool {
	errorKeywords := []string{
		"Ошибка получени",
		"ошибка получени",
		"Error getting",
		"error getting",
		"Failed to get",
		"failed to get",
	}
	
	legalNameLower := strings.ToLower(legalName)
	countryLower := strings.ToLower(country)
	
	for _, keyword := range errorKeywords {
		keywordLower := strings.ToLower(keyword)
		if strings.Contains(legalNameLower, keywordLower) || 
		   strings.Contains(countryLower, keywordLower) {
			return true
		}
	}
	
	return false
}

func (r *DomainRepository) Create(domain *model.Domain) error {
	// Проверяем на ошибочные данные перед созданием
	if r.isErrorData(domain.LegalName, domain.Country) {
		// Просто пропускаем запись без ошибки
		return nil
	}
	
	return r.db.Create(domain).Error
}

func (r *DomainRepository) Save(domain *model.Domain) error {
	// Проверяем на ошибочные данные перед сохранением
	if r.isErrorData(domain.LegalName, domain.Country) {
		// Просто пропускаем запись без ошибки
		return nil
	}
	
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

func (r *DomainRepository) GetByDomain(domain string) (*model.Domain, error) {
	var d model.Domain
	if err := r.db.Where("domain = ?", domain).First(&d).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &d, nil
}

func (r *DomainRepository) TruncateAll() error {
	return r.db.Exec("TRUNCATE TABLE domains RESTART IDENTITY").Error
}

// CreateBatch создает несколько записей сразу, пропуская ошибочные
func (r *DomainRepository) CreateBatch(domains []model.Domain) (int, error) {
	validDomains := make([]model.Domain, 0)
	
	for _, domain := range domains {
		if !r.isErrorData(domain.LegalName, domain.Country) {
			validDomains = append(validDomains, domain)
		}
	}
	
	if len(validDomains) == 0 {
		return 0, errors.New("нет корректных данных для сохранения")
	}
	
	err := r.db.Create(&validDomains).Error
	return len(validDomains), err
}