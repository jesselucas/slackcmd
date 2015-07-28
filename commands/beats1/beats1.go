package beats1

import (
	"fmt"
	"github.com/jesselucas/slackcmd/slack"
)

type Command struct {
}

func (cmd Command) Request(sc *slack.SlashCommand) (*slack.CommandPayload, error) {
	// create payload
	cp := &slack.CommandPayload{
		Channel:       fmt.Sprintf("@%v", sc.UserName),
		Username:      "Beats1",
		Emoji:         ":metal:",
		SlashResponse: false,
		SendPayload:   true,
	}

	cp.Text = "Beats1"

	// Read the latest tweet
	// https://api.twitter.com/1.1/statuses/user_timeline.json
	// screen_name=beats1plays&count=1& trim_user=true

	// read credentials from environment variables
	// consumerKey := os.Getenv("TWITTER_CONSUMER_KEY")
	// consumerSecret := os.Getenv("TWITTER_CONSUMER_SECRET")
	// accessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	// accessTokenSecret := os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")
	// if consumerKey == "" || consumerSecret == "" || accessToken == "" || accessTokenSecret == "" {
	// 	panic("Missing required environment variable")
	// }

	return cp, nil
}
