package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jesselucas/slackcmd/beats1"
	"github.com/jesselucas/slackcmd/slack"
	"github.com/jesselucas/slackcmd/trello"
	"log"
	"net/http"
	"net/url"
	"os"
)

func createSlashCommand(w http.ResponseWriter, r *http.Request) *slack.SlashCommand {
	var v url.Values

	switch r.Method {
	case "POST":
		r.ParseForm()
		v = r.Form
	case "GET":
		v = r.URL.Query()
	}

	sc := &slack.SlashCommand{
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
	sc := createSlashCommand(w, r)

	// check url to see what command
	cmdURL := r.URL.Path[len("/cmd/"):]

	fmt.Println("command", cmdURL)

	// Each command implements the handler interface ServeHTTP(ResponseWriter, *Request)
	var command slack.Command

	// Add commands here
	switch cmdURL {
	case "trello":
		command = trello.Trello{}
	case "beats1":
		command = beats1.Beats1{}
	}

	fmt.Println("slash command:", sc.Text)

	// command request returns payload
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
