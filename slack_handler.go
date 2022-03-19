package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/st-phuongvu/st-aws-slack-bot/model"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// handleEventMessage will take an event and handle it properly based on the type of event
func (b *Bot) HandleEventMessage(event slackevents.EventsAPIEvent) (string, error) {
	switch event.Type {
	case slackevents.CallbackEvent:
		innerEvent := event.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			client := b.SlackClient.Client
			user, err := client.GetUserInfo(ev.User)
			if err != nil {
				return ev.Channel, err
			}
			text := strings.ToLower(ev.Text)

			// exclude bot mention part
			params := [3]string{}
			copy(params[:], regexp.MustCompile("\\s[\\d\\w-,]+").FindAllString(text, 3))
			return ev.Channel, b.HandleEventRequest(*user, ev.Channel, params)
		}
	default:
		return "", errors.New("unsupported event type")
	}
	return "", nil
}

func (b *Bot) HandleInteractiveMessage(interaction slack.InteractionCallback) (string, error) {
	switch interaction.Type {
	case slack.InteractionTypeBlockActions:
		action := interaction.ActionCallback
		if len(action.BlockActions) > 0 {
			triggerAction := action.BlockActions[0]
			switch triggerAction.ActionID {
			case BLOCK_ACTION_ID_NEW_RESOURCE:
				state := make(map[string]interface{})
				err := json.Unmarshal([]byte(interaction.RawState), &state)
				if err != nil {
					return interaction.Channel.ID, errors.New("ERROR WHEN UNMARSHAL SLACK RESPONSE")
				}
				var resourceID, accessKey, secretKey string
				for _, v1 := range state {
					for _, v2 := range v1.(map[string]interface{}) {
						for k, v3 := range v2.(map[string]interface{}) {
							b, err := json.Marshal(v3.(map[string]interface{}))
							if err != nil {
								return interaction.Channel.ID, errors.New("ERROR WHEN DECODE STATE")
							}
							data := map[string]string{}
							err = json.Unmarshal(b, &data)
							if err != nil {
								return interaction.Channel.ID, errors.New("ERROR WHEN DECODE STATE")
							}
							switch k {
							case ACTION_ID_RESOURCE:
								resourceID = data["value"]
							case ACTION_ID_ACCESS_KEY:
								accessKey = data["value"]
							case ACTION_ID_SECRET_KEY:
								secretKey = data["value"]
							}
						}
					}
				}
				err = b.HandleNewResourceRequest(
					interaction.Channel.ID,
					&model.AWSResource{
						ID:        resourceID,
						AccessKey: accessKey,
						SecretKey: secretKey,
					},
				)
				if err != nil {
					return interaction.Channel.ID, err
				}
				return interaction.Channel.ID, nil
			case BLOCK_ACTION_ID_REBOOT_EC2:
				user := interaction.User
				channel := interaction.Channel.ID
				params := [3]string{
					"ec2",
					EC2_REBOOT_ACTION,
					triggerAction.Value,
				}
				return interaction.Channel.ID, b.HandleEventRequest(user, channel, params)
			default:
				return interaction.Channel.ID, errors.New("not a valid action")
			}
		}
		return interaction.Channel.ID, errors.New("not a valid action")
	default:
		return interaction.Channel.ID, errors.New("unsupported event type")
	}
}

func (b *Bot) HandleNewResourceRequest(channel string, resource *model.AWSResource) error {
	err := b.AWSResourceRepository.Insert(resource)
	if err != nil {
		return err
	}
	reply := fmt.Sprintf("SUCCESSFULLY ADD INSTANCE %v", resource.ID)
	_, _, err = b.SlackClient.Client.PostMessage(channel, slack.MsgOptionText(reply, false))
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}
	return nil
}

func (b *Bot) HandleEventRequest(user slack.User, channel string, params [3]string) error {
	if params[0] == "" {
		return errors.New("No params")
	}
	switch strings.TrimSpace(params[0]) {
	case "ec2":
		resource, err := b.AWSResourceRepository.Get(params[2])
		if err != nil {
			return err
		}
		resp, err := NewEC2Handler(
			user.ID,
			strings.TrimSpace(params[1]),
			resource,
		)
		if err != nil {
			return err
		}
		succeed := []string{}
		failed := []string{}
		for id, ok := range resp.InstancesState {
			if ok {
				succeed = append(succeed, id)
			} else {
				failed = append(failed, id)
			}
		}
		reply := fmt.Sprintf("RESOURCE: %v\nACTION: %v\nSUCCESS INSTANCE: %v\nFAILED INSTANCE: %v",
			strings.TrimSpace(params[0]),
			strings.TrimSpace(params[1]),
			strings.Join(succeed, ","),
			strings.Join(failed, ","),
		)
		_, _, err = b.SlackClient.Client.PostMessage(channel, slack.MsgOptionText(reply, false))
		if err != nil {
			return fmt.Errorf("failed to post message: %w", err)
		}
	case "new":
		_, _, err := b.SlackClient.Client.PostMessage(channel, slack.MsgOptionBlocks(BuildBlockNewEvent()...))
		if err != nil {
			return fmt.Errorf("failed to post message: %w", err)
		}
	case "list":
		resources, err := b.AWSResourceRepository.Find(map[string]interface{}{})
		if err != nil {
			return err
		}
		ids := []string{}
		for _, res := range resources {
			ids = append(ids, res.ID)
		}
		_, _, err = b.SlackClient.Client.PostMessage(channel, slack.MsgOptionBlocks(BuildBlockListEvent(ids)...))
		if err != nil {
			return fmt.Errorf("failed to post message: %w", err)
		}
	default:
		return errors.New("Unsupported resource type")
	}

	return nil
}
