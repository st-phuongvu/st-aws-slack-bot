package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func (b *Bot) initSlack() (context.CancelFunc, error) {
	client := slack.New(
		os.Getenv("SLACK_BOT_TOKEN"),
		slack.OptionDebug(true),
		slack.OptionAppLevelToken(os.Getenv("SLACK_BOT_SOCKET_MODE_TOKEN")),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
	)

	socketClient := socketmode.New(
		client,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	// Create a context that can be used to cancel goroutine
	ctx, cancel := context.WithCancel(context.Background())
	// Make this cancel called properly in a real program , graceful shutdown etc

	go func(ctx context.Context, client *slack.Client, socketClient *socketmode.Client) {
		// Create a for loop that selects either the context cancellation or the events incomming
		for {
			select {
			// inscase context cancel is called exit the goroutine
			case <-ctx.Done():
				log.Println("Shutting down socketmode listener")
				return
			case event := <-socketClient.Events:
				// We have a new Events, let's type switch the event
				// Add more use cases here if you want to listen to other events.
				switch event.Type {
				// handle EventAPI events
				case socketmode.EventTypeEventsAPI:
					// The Event sent on the channel is not the same as the EventAPI events so we need to type cast it
					eventsAPIEvent, ok := event.Data.(slackevents.EventsAPIEvent)
					if !ok {
						log.Printf("Could not type cast the event to the EventsAPIEvent: %v\n", event)
						continue
					}
					// We need to send an Acknowledge to the slack server
					socketClient.Ack(*event.Request)
					// Now we have an Events API event, but this event type can in turn be many types, so we actually need another type switch
					channel, err := b.HandleEventMessage(eventsAPIEvent)
					if err != nil {
						if channel != "" {
							client.PostMessage(channel, slack.MsgOptionText(err.Error(), false))
						} else {
							log.Println(err.Error())
						}
					}
				case socketmode.EventTypeInteractive:
					eventsInteractive, ok := event.Data.(slack.InteractionCallback)
					if !ok {
						log.Printf("Could not type cast the event to the EventsAPIEvent: %v\n", event)
						continue
					}
					channel, err := b.HandleInteractiveMessage(eventsInteractive)
					if err != nil {
						if channel != "" {
							client.PostMessage(channel, slack.MsgOptionText(err.Error(), false))
						} else {
							log.Println(err.Error())
						}
						continue
					}
					// Delete message after finish
					data := map[string]string{
						"delete_original": "true",
					}
					json_data, err := json.Marshal(data)
					http.DefaultClient.Post(
						eventsInteractive.ResponseURL,
						"application/json",
						bytes.NewBuffer(json_data),
					)
				}
			}
		}
	}(ctx, client, socketClient)

	b.SlackClient = socketClient

	return cancel, nil
}
