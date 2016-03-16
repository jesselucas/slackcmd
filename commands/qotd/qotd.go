package qotd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/forestgiant/go-simpletime"
	"github.com/jesselucas/slackcmd/slack"
	"github.com/jesselucas/validator"
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
	if !validator.IsURL(url) {
		return nil, errors.New("QOTD_URL is not a valid URL")
	}

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

	// Get todays index
	index := getTodaysIndex(uint(len(questions)))

	cp.Text = "QOTD: " + questions[index]

	return cp, nil
}

// getTodaysIndex subjects the startDate by today's date to get the
// difference in days and then will modulate based on the length
// of all the indexes
func getTodaysIndex(length uint) uint {
	startDate := time.Date(2016, 3, 16, 0, 0, 0, 0, time.UTC)

	// Find out day offset of today from the startDate
	offsetDuration := simpletime.NewSimpleTime(startDate).Since(time.Now())
	offsetDays := offsetDuration.Days()
	fmt.Println(offsetDays)

	return uint(offsetDays) % length
}
