package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"parser/internal/adapter/repository"
	"parser/internal/model"
)

type DomainService struct {
	repo *repository.DomainRepository
}

func NewDomainService(repo *repository.DomainRepository) *DomainService {
	return &DomainService{repo: repo}
}

type ParsedResult struct {
	LegalName string `json:"legal_name"`
	Country   string `json:"country"`
	Error     string `json:"error,omitempty"`
}

func (s *DomainService) ParseAndSave(domain string) (*model.Domain, error) {
	existing, err := s.repo.GetByDomain(domain)
	if err == nil && existing != nil {
	return existing, nil
	}

	url := fmt.Sprintf("http://parser:3001/parse?domain=%s", domain)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http request to parser failed: %w", err)
	}
	defer resp.Body.Close()

	var result ParsedResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode parser response: %w", err)
	}
	if result.Error != "" {
		return nil, fmt.Errorf("parser returned error: %s", result.Error)
	}

	entity := &model.Domain{
		Domain:    domain,
		LegalName: result.LegalName,
		Country:   result.Country,
	}


	if err := s.repo.Save(entity); err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *DomainService) StreamParseBatchAndSave(domains []string) ([]model.Domain, error) {
	if len(domains) == 0 {
		return nil, errors.New("no domains provided")
	}

	// JSON тело запроса
	payload := map[string][]string{"domains": domains}
	body, _ := json.Marshal(payload)

	// Запрос к Node.js серверу
	resp, err := http.Post("http://parser:3001/parse-batch", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to request parser: %w", err)
	}
	defer resp.Body.Close()

	// Чтение JSON массива поэлементно
	decoder := json.NewDecoder(resp.Body)

	// Ждём открывающую скобку массива
	tok, err := decoder.Token()
	if err != nil || tok != json.Delim('[') {
		return nil, fmt.Errorf("expected [ token, got %v", tok)
	}

	var saved []model.Domain

	for decoder.More() {
		var item model.ParsedBatchItem
		if err := decoder.Decode(&item); err != nil {
			return nil, fmt.Errorf("decode stream item failed: %w", err)
		}
		if item.Error != "" {
			log.Printf("⚠️ error for domain %s: %s", item.Domain, item.Error)
			continue
		}

		// проверка на дубликат
		existing, err := s.repo.GetByDomain(item.Domain)
		if err == nil && existing != nil {
			saved = append(saved, *existing)
			continue
		}

		entity := model.Domain{
			Domain:    item.Domain,
			LegalName: item.LegalName,
			Country:   item.Country,
		}

		if err := s.repo.Save(&entity); err == nil {
			saved = append(saved, entity)
		}
	}

	return saved, nil
}


func (s *DomainService) ParseBatchAndSave(domains []string) ([]model.Domain, error) {
	if len(domains) == 0 {
		return nil, errors.New("no domains provided")
	}

	payload := map[string][]string{"domains": domains}
	body, _ := json.Marshal(payload)

	resp, err := http.Post("http://parser:3001/parse-batch", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to request parser: %w", err)
	}
	defer resp.Body.Close()

	var parsed []model.ParsedBatchItem
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("failed to decode parser response: %w", err)
	}

	var saved []model.Domain
	for _, item := range parsed {
		if item.Error != "" {
			continue // пропускаем невалидные
		}

		// проверка на дубль
		existing, err := s.repo.GetByDomain(item.Domain)
		if err == nil && existing != nil {
			saved = append(saved, *existing)
			continue
		}

		entity := model.Domain{
			Domain:    item.Domain,
			LegalName: item.LegalName,
			Country:   item.Country,
		}

		if err := s.repo.Save(&entity); err == nil {
			saved = append(saved, entity)
		}
	}

	return saved, nil
}

func (s *DomainService) Save(domain *model.Domain) error {
	return s.repo.Save(domain)
}

//GetAll
func (s *DomainService) GetAll() ([]model.Domain, error) {
	return s.repo.GetAll()
}

func (s *DomainService) DeleteByID(id uint) error {
	return s.repo.DeleteByID(id)
}

func (s *DomainService) GetPaginated(domain string, limit, offset int) ([]model.Domain, int64, error) {
	return s.repo.GetPaginated(domain, limit, offset)
}

func (s *DomainService) CleanData() error {
	return s.repo.TruncateAll()
}

