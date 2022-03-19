package main

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/st-phuongvu/st-aws-slack-bot/model"
)

func (b *Bot) initDB() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSL_MODE")
	pass := os.Getenv("DB_PASS")
	dsn := fmt.Sprintf("host=%v port=%v user=%v dbname=%v sslmode=%v password=%v", host, port, user, dbName, sslmode, pass)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&model.AWSResource{}, &model.Management{})

	return db, nil
}
