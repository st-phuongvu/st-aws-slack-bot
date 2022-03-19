package main

import (
	"log"

	"github.com/st-phuongvu/st-aws-slack-bot/repository"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	b := &Bot{}
	db, err := b.initDB()
	if err != nil {
		log.Fatal(err)
	}
	repo := &repository.AWSResourceRepository{
		Db: db,
	}
	b.AWSResourceRepository = repo

	cancel, err := b.initSlack()
	defer cancel()
	if err != nil {
		log.Fatal(err)
	}

	b.SlackClient.Run()
}
