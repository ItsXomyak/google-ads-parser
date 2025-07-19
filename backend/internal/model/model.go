package model

import "time"

type Domain struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Domain     string    `gorm:"uniqueIndex" json:"domain"`
	LegalName  string    `json:"legal_name"`
	Country    string    `json:"country"`
	CreatedAt  time.Time `json:"created_at"`
}

type ParsedBatchItem struct {
	Domain     string `json:"domain"`
	LegalName  string `json:"legal_name,omitempty"`
	Country    string `json:"country,omitempty"`
	Error      string `json:"error,omitempty"`
}
