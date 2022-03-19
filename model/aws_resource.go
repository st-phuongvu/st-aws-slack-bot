package model

import "time"

type AWSResource struct {
	ID        string
	AccessKey string
	SecretKey string
	CreatedAt time.Time `gorm:"autoUpdateTime"`
	UpdatedAt time.Time `gorm:"autoCreateTime"`
}
