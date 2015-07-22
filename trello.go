package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type Trello struct {
}

type board struct {
	Name string
	Id   string
}

type list struct {
	Name    string
	Id      string
	IdBoard string
	Cards   []card
}

type card struct {
	Name string
	Id   string
}

func (t Trello) Request(sc *SlashCommand) (*CommandPayload, error) {

	c := strings.Fields(sc.Text)

	// replace all underscores "_" with spaces " " in commands
	commands := make([]string, len(c))
	copy(commands, c)

	for i := 0; i < len(c); i++ {
		c[i] = strings.Replace(c[i], "_", " ", -1)
	}

	cp := &CommandPayload{
		Username: "FGBot",
		Emoji:    ":fg:",
	}

	// construct url for Trello
	url := fmt.Sprintf(
		"https://api.trello.com/1/organizations/%v/boards/?fields=name&filter=open&key=%v&token=%v",
		"forestgiant",
		os.Getenv("TRELLO_KEY"),
		os.Getenv("TRELLO_TOKEN"),
	)

	res, err := http.Get(url)

	body, _ := ioutil.ReadAll(res.Body)

	// found boards return if only sent one command
	var boards []board
	json.Unmarshal(body, &boards)

	var responseString string

	// if the command is blank return possible commands(boards)
	if len(c) == 0 {

		// iterate over boards and create string to send
		for _, board := range boards {
			responseString += fmt.Sprintf(
				"* <https://trello.com/b/%v|%v> \n",
				board.Id,
				board.Name,
			)
		}

		fmt.Println("response string:", responseString)
		cp.Text = responseString

		return cp, nil
	}

	if err != nil {
		return nil, err
	}

	return cp, nil
}

// IsUrl test if the rxURL regular expression matches a string
func IsURL(s string) bool {
	rxURL := regexp.MustCompile(`\b(([\w-]+://?|www[.])[^\s()<>]+(?:\([\w\d]+\)|([^[:punct:]\s]|/)))`)

	if s == "" || len(s) >= 2083 || len(s) <= 10 || strings.HasPrefix(s, ".") {
		return false
	}
	u, err := url.Parse(s)
	if err != nil {
		return false
	}
	if strings.HasPrefix(u.Host, ".") {
		return false
	}
	return rxURL.MatchString(s)
}
