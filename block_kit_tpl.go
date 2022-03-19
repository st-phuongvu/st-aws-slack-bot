package main

import (
	"fmt"

	"github.com/slack-go/slack"
)

const (
	// SLACK BOT TYPE FOR BUILDING BLOCK
	INPUT_TYPE            = "input"
	SECTION_TYPE          = "section"
	MARKDOWN_TEXT_TYPE    = "mrkdwn"
	PLAIN_TEXT_INPUT_TYPE = "plain_text_input"
	PLAIN_TEXT_TYPE       = "plain_text"
	BUTTON_TYPE           = "button"

	// DEFINE ACTION ID FOR USING ON HANDLE EVENT
	ACTION_ID_RESOURCE           = "resource_id"
	ACTION_ID_ACCESS_KEY         = "access_key"
	ACTION_ID_SECRET_KEY         = "secret_key"
	BLOCK_ACTION_ID_NEW_RESOURCE = "submit-new-resource"
	BLOCK_ACTION_ID_REBOOT_EC2   = "reboot-ec2"
)

type BlockResponse struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func BuildBlockListEvent(ids []string) []slack.Block {
	blocks := []slack.Block{}
	for _, id := range ids {
		block := slack.SectionBlock{
			Type: SECTION_TYPE,
			Text: &slack.TextBlockObject{
				Type: MARKDOWN_TEXT_TYPE,
				Text: fmt.Sprintf("INSTANCE ID: %v", id),
			},
			Accessory: &slack.Accessory{
				ButtonElement: &slack.ButtonBlockElement{
					Type: BUTTON_TYPE,
					Text: &slack.TextBlockObject{
						Type: PLAIN_TEXT_TYPE,
						Text: "Reboot",
					},
					Value:    id,
					ActionID: BLOCK_ACTION_ID_REBOOT_EC2,
				},
			},
		}
		blocks = append(blocks, block)
	}
	return blocks
}

func BuildBlockNewEvent() []slack.Block {
	idBlock := slack.InputBlock{
		Type: INPUT_TYPE,
		Element: slack.PlainTextInputBlockElement{
			Type:     PLAIN_TEXT_INPUT_TYPE,
			ActionID: ACTION_ID_RESOURCE,
		},
		Label: &slack.TextBlockObject{
			Type: PLAIN_TEXT_TYPE,
			Text: "Resource ID",
		},
	}

	accessKeyBlock := slack.InputBlock{
		Type: INPUT_TYPE,
		Element: slack.PlainTextInputBlockElement{
			Type:     PLAIN_TEXT_INPUT_TYPE,
			ActionID: ACTION_ID_ACCESS_KEY,
		},
		Label: &slack.TextBlockObject{
			Type: PLAIN_TEXT_TYPE,
			Text: "Access Key",
		},
	}

	secretKeyBlock := slack.InputBlock{
		Type: INPUT_TYPE,
		Element: slack.PlainTextInputBlockElement{
			Type:     PLAIN_TEXT_INPUT_TYPE,
			ActionID: ACTION_ID_SECRET_KEY,
		},
		Label: &slack.TextBlockObject{
			Type: PLAIN_TEXT_TYPE,
			Text: "Secret Key",
		},
	}

	submitButton := slack.SectionBlock{
		Type: SECTION_TYPE,
		Text: &slack.TextBlockObject{
			Type: MARKDOWN_TEXT_TYPE,
			Text: "Recheck before submit resource",
		},
		Accessory: &slack.Accessory{
			ButtonElement: &slack.ButtonBlockElement{
				Type: BUTTON_TYPE,
				Text: &slack.TextBlockObject{
					Type: PLAIN_TEXT_TYPE,
					Text: "Submit",
				},
				Value:    BLOCK_ACTION_ID_NEW_RESOURCE,
				ActionID: BLOCK_ACTION_ID_NEW_RESOURCE,
			},
		},
	}

	return []slack.Block{
		idBlock, accessKeyBlock, secretKeyBlock, submitButton,
	}
}
