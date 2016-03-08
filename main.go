package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/jesselucas/slackcmd/commands/beats1"
	"github.com/jesselucas/slackcmd/commands/calendar"
	"github.com/jesselucas/slackcmd/commands/qotd"
	"github.com/jesselucas/slackcmd/commands/trello"
	"github.com/jesselucas/slackcmd/slack"
)

// struct used to store environment variables from config.json
type env struct {
	Key   string
	Value string
}

func setEnvFromJSON(configPath string) {
	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Println("config.json not found. Using os environment variables.")
		return
	}

	var envVars []env
	json.Unmarshal(configFile, &envVars)

	// set environment variables
	for _, env := range envVars {
		// fmt.Println(env)
		os.Setenv(env.Key, env.Value)
	}

}

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

	// interface reference for slack Commands
	var cmd slack.Command

	// Create FlagSet to store flags
	fs := &slack.FlagSet{}

	// Add commands here
	switch sc.Command {
	case "/fg":
		cmd = &trello.Command{}
		fs.Usage = "/fg help: FG Trello access"
	case "/beats1":
		cmd = &beats1.Command{}
		fs.Usage = "/beats1 help: Song currently playing on Beats1"
	case "/conference":
		cmd = &calendar.Command{}
		fs.Usage = "/conference help: Schedule for FG Conference room"
	case "/qotd":
		cmd = &qotd.Command{}
		fs.Usage = "/qotd help: Sends the Question of the Day"
	default:
		err := errors.New("No Command found")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	fmt.Println("slash command:", sc.Text)

	// parse out flags
	parsedCommands, flags := slack.SeparateFlags(sc.Text)
	sc.Text = parsedCommands

	// command request returns payload
	cp, err := cmd.Request(sc)

	if cp == nil {
		err := errors.New("Unauthorized")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	if err != nil {
		err := errors.New("Unauthorized")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// Set Flags for Commands
	slack.SetFlag(fs, "channel", "c", "Sends the response to the current channel", func() {
		cp.Channel = fmt.Sprintf("#%v", sc.ChannelName)
		cp.SendPayload = true
		cp.SlashResponse = false
	})

	slack.SetFlag(fs, "private", "p", "Sends a private message with the response", func() {
		cp.SendPayload = true
		cp.SlashResponse = false
	})

	sc.Text = parsedCommands

	// TODO: Move ParseFlags call so it happen before cmd.Request is called
	help, response := slack.ParseFlags(fs, flags)
	if help == true {
		cp.Text = response
	}

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

func main() {
	// setup environment variables if a config json exist
	setEnvFromJSON("config.json")

	// url setup. FIX make more generic
	var url string
	if os.Getenv("PORT") != "" {
		url = ":" + os.Getenv("PORT")
	} else {
		url = "localhost:8080"
	}

	// vs := validateSlackToken(http.HandlerFunc(commandHandler), slackAPIKey)
	http.HandleFunc("/cmd/", commandHandler)
	http.HandleFunc("/cmd", commandHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Go away!")
	})
	log.Fatal(http.ListenAndServe(url, nil))
}
