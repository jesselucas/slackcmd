package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

// interface for commands
type Command interface {
	Request(sc *SlashCommand) (*CommandPayload, error)
}

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type Attachment struct {
	Fallback  string  `json:"fallback"`
	Title     string  `json:"title"`
	TitleLink string  `json:"title_link"`
	Text      string  `json:"text"`
	Pretext   string  `json:"pretext"`
	Color     string  `json:"color"`
	Fields    []Field `json:"fields"`
}

// CommandPayload
type CommandPayload struct {
	Channel       string       `json:"channel"`
	Username      string       `json:"username"`
	Emoji         string       `json:"icon_emoji"`
	EmojiURL      string       `json:"icon_url"`
	Text          string       `json:"text"`
	Attachments   []Attachment `json:"attachments"`
	UnfurlMedia   bool         `json:"unfurl_media"`
	UnfurlLinks   bool         `json:"unfurl_links"`
	SlashResponse bool
	SendPayload   bool
}

// struct to hold params sent from slacks slash command
type SlashCommand struct {
	Token       string
	TeamId      string
	TeamDomain  string
	ChannelId   string
	ChannelName string
	UserId      string
	UserName    string
	Command     string
	Text        string
	Hook        string
}

func createSlashCommand(w http.ResponseWriter, r *http.Request) *SlashCommand {
	var v url.Values

	switch r.Method {
	case "POST":
		r.ParseForm()
		v = r.Form
	case "GET":
		v = r.URL.Query()
	}

	sc := &SlashCommand{
		Token:       v.Get("token"),
		TeamId:      v.Get("team_id"),
		TeamDomain:  v.Get("team_domain"),
		ChannelId:   v.Get("channel_id"),
		ChannelName: v.Get("channel_name"),
		UserId:      v.Get("user_id"),
		UserName:    v.Get("user_name"),
		Command:     v.Get("command"),
		Text:        v.Get("text"),
		Hook:        v.Get("hook"),
	}

	return sc
}

func commandHandler(w http.ResponseWriter, r *http.Request) {
	slackAPIKey := os.Getenv("SLACK_KEY")

	sc := createSlashCommand(w, r)
	// Verify the request is coming from Slack
	if sc.Token != slackAPIKey {
		err := errors.New("Unauthorized")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// check url to see what command
	cmdURL := r.URL.Path[len("/cmd/"):]

	fmt.Println("command", cmdURL)

	// Each command implements the handler interface ServeHTTP(ResponseWriter, *Request)
	var command Command

	switch cmdURL {
	case "trello":
		command = Trello{}
	case "beats1":
		command = Beats1{}
	}

	fmt.Println("slash command:", sc.Text)

	// command request return payload
	cp, err := command.Request(sc)

	if cp == nil {
		err := errors.New("Unauthorized")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	if err != nil {
		err := errors.New("Unauthorized")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	} else {

		// check if the command wants to send a slash command response
		if cp.SlashResponse {
			w.Write([]byte(cp.Text))
		}

		// don't send payload if hook URL isn't passed
		if sc.Hook != "" && cp.SendPayload == true {
			cpJSON, err := json.Marshal(cp)
			if err != nil {
				err := errors.New("Unauthorized")
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}

			cpJSONString := string(cpJSON[:])

			// Make the request to the Slack API.
			http.PostForm(sc.Hook, url.Values{"payload": {cpJSONString}})
		}
	}
}

func main() {
	// url setup. FIX make more generic
	var url string
	if os.Getenv("PORT") != "" {
		url = ":" + os.Getenv("PORT")
	} else {
		url = "localhost:8080"
	}

	// vs := validateSlackToken(http.HandlerFunc(commandHandler), slackAPIKey)
	http.HandleFunc("/cmd/", commandHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Go away!")
	})
	log.Fatal(http.ListenAndServe(url, nil))
}
