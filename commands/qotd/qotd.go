package qotd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"

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

	// url for QOTD YAML
	url := os.Getenv("QOTD_URL")

	res, err := http.Get(url)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal YAML
	var questions []string
	err = yaml.Unmarshal(body, &questions)
	if err != nil {
		return nil, err
	}

	cp.Text = questions[rand.Intn(len(questions))]

	return cp, nil
}
