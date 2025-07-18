package model

import "time"

type Domain struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Domain     string    `gorm:"uniqueIndex" json:"domain"`
	LegalName  string    `json:"legal_name"`
	Country    string    `json:"country"`
	CreatedAt  time.Time `json:"created_at"`
}
