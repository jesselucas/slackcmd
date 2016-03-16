package qotd

import (
	"fmt"
	"os"
	"testing"

	"github.com/jesselucas/slackcmd/slack"
)

var cmd *Command
var sc *slack.SlashCommand

func init() {
	// Create Command
	cmd = new(Command)

	sc = &slack.SlashCommand{
		Token:       "Js7gTRur9cWBjXnWdYfm2XXy",
		TeamId:      "T0001",
		TeamDomain:  "example",
		ChannelId:   "C2147483705",
		ChannelName: "test",
		UserId:      "U2147483697",
		UserName:    "Steve",
		Command:     "/qotd",
		Text:        "",
		Hook:        "https://hooks.slack.com/commands/1234/5678",
	}

	// setup environment variables
	// Make sure you have QOTD_URL set as an environment variable
	// os.Setenv("QOTD_URL", "http://urltoquestions.yaml")
	os.Setenv("SLACK_KEY_QOTD", "Js7gTRur9cWBjXnWdYfm2XXy")
}

func TestRequest(t *testing.T) {
	cp, err := cmd.Request(sc)

	if cp == nil {
		t.Error("payload is nil")
	}

	if err != nil {
		t.Error("Request error:", err)
	}

	fmt.Println("payload", cp)
}
