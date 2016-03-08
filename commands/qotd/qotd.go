package qotd

import (
	"errors"
	"fmt"
	"os"

	"github.com/jesselucas/slackcmd/slack"
)

// Command struct is only defined to add the Request method
type Command struct {
}

// Request is used to send back to slackcmd
func (cmd *Command) Request(sc *slack.SlashCommand) (*slack.CommandPayload, error) {
	// read credentials from environment variables
	slackAPIKey := os.Getenv("SLACK_KEY_QOTD")

	// Verify the request is coming from Slack
	if sc.Token != slackAPIKey {
		err := errors.New("Unauthorized Slack")
		return nil, err
	}

	// create payload
	cp := &slack.CommandPayload{
		Channel:       fmt.Sprintf("@%v", sc.UserName),
		Username:      "QOTD",
		Emoji:         ":question:",
		SlashResponse: true,
		SendPayload:   false,
	}

	cp.Text = "What is your favorite star?"

	return cp, nil
}
