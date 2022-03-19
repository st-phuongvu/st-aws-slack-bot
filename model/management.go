package model

import "time"

type Management struct {
	UserID        string
	Email         string
	AWSResourceID string
	CreatedAt     time.Time `gorm:"autoUpdateTime"`
	UpdatedAt     time.Time `gorm:"autoCreateTime"`
}
