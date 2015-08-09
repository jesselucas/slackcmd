package beats1

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dghubble/oauth1"
	"github.com/jesselucas/slackcmd/slack"
)

type Command struct {
}

func (cmd Command) Request(sc *slack.SlashCommand) (*slack.CommandPayload, error) {
	// read credentials from environment variables
	slackAPIKey := os.Getenv("SLACK_KEY_BEATS1")
	consumerKey := os.Getenv("TWITTER_CONSUMER_KEY")
	consumerSecret := os.Getenv("TWITTER_CONSUMER_SECRET")
	accessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret := os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")
	if slackAPIKey == "" || consumerKey == "" || consumerSecret == "" || accessToken == "" || accessTokenSecret == "" {
		panic("Missing required environment variable")
	}

	// Verify the request is coming from Slack
	if sc.Token != slackAPIKey {
		err := errors.New("Unauthorized Slack")
		return nil, err
	}

	// create payload
	cp := &slack.CommandPayload{
		Channel:       fmt.Sprintf("@%v", sc.UserName),
		Username:      "Beats1",
		Emoji:         ":metal:",
		SlashResponse: true,
		SendPayload:   false,
	}

	cp.Text = "Beats1"

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessTokenSecret)

	// httpClient will automatically authorize http.Request's
	httpClient := config.Client(token)

	url := fmt.Sprintf(
		"https://api.twitter.com/1.1/statuses/user_timeline.json?screen_name=%v&count=%v&trim_user=%v",
		"beats1plays",
		"1",
		"true",
	)
	res, err := httpClient.Get(url)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var tweets []tweet
	json.Unmarshal(body, &tweets)

	var responseString string
	for _, t := range tweets {
		responseString += fmt.Sprint(t)
	}

	cp.Text = string(responseString)

	return cp, nil
}

type tweet struct {
	Text     string `json:"text"`
	Entities entity `json:"entities"`
}

func (t tweet) String() string {

	// remove hashtag
	text := strings.Replace(t.Text, "#beats1", "", -1)

	// remove Urls
	urls := t.Entities.Urls
	var expandedURL string
	for _, u := range urls {
		urlString := u.Url
		text = strings.Replace(text, urlString, "", -1)
		expandedURL += u.ExpandedURL
	}

	// remove media
	media := t.Entities.Media
	for _, m := range media {
		mediaURL := m.Url
		text = strings.Replace(text, mediaURL, "", -1)
	}

	// trime white space
	text = strings.TrimSpace(text)

	// if there is an expandedURL
	if expandedURL != "" {
		return fmt.Sprintf(
			"<%v|%v>",
			expandedURL,
			text,
		)
	}

	return fmt.Sprintf(
		"%v",
		text,
	)

}

type entity struct {
	Urls  []url   `json:"urls"`
	Media []media `json:"media"`
}

type url struct {
	Url         string `json:"url"`
	ExpandedURL string `json:"expanded_url"`
}

type media struct {
	Url string `json:"url"`
}
