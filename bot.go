package main

import (
	"github.com/st-phuongvu/st-aws-slack-bot/repository"

	"github.com/slack-go/slack/socketmode"
)

type Bot struct {
	SlackClient           *socketmode.Client
	AWSResourceRepository repository.AWSResourceRepositoryInterface
}
