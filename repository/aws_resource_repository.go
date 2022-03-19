package repository

import (
	"github.com/st-phuongvu/st-aws-slack-bot/model"

	"gorm.io/gorm"
)

type AWSResourceRepositoryInterface interface {
	Get(id string) (*model.AWSResource, error)
	Insert(resource *model.AWSResource) error
	Find(condition map[string]interface{}) ([]model.AWSResource, error)
}

type AWSResourceRepository struct {
	Db *gorm.DB
}

func (r *AWSResourceRepository) Find(condition map[string]interface{}) ([]model.AWSResource, error) {
	resources := []model.AWSResource{}
	tx := r.Db.Where(condition).Find(&resources)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return resources, nil
}

func (r *AWSResourceRepository) Insert(resource *model.AWSResource) error {
	return r.Db.Create(&resource).Error
}

func (r *AWSResourceRepository) Get(id string) (*model.AWSResource, error) {
	resource := &model.AWSResource{}
	tx := r.Db.First(&resource)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return resource, nil
}
