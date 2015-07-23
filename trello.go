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

func (b board) String() string {
	return fmt.Sprintf(
		"• <https://trello.com/b/%v|%v> \n",
		b.Id,
		strings.Replace(b.Name, " ", "_", -1),
	)
}

type list struct {
	Name    string
	Id      string
	IdBoard string
	Cards   []card
}

func (l list) String() string {
	return fmt.Sprintf(
		"• <https://trello.com/b/%v|%v> \n",
		l.IdBoard,
		strings.Replace(l.Name, " ", "_", -1),
	)
}

type card struct {
	Name string
	Id   string
	URL  bool
}

func (c card) String() string {
	var url string

	if c.URL {
		url = c.Name
	} else {
		url = fmt.Sprintf("https://trello.com/c/%v", c.Id)
	}

	return fmt.Sprintf(
		"• <%v|%v> \n",
		url,
		c.Name,
	)
}

func (t Trello) Request(sc *SlashCommand) (*CommandPayload, error) {

	// create payload
	cp := &CommandPayload{
		Channel:       fmt.Sprintf("@%v", sc.UserName),
		Username:      "FG Bot",
		Emoji:         ":fgdot:",
		SlashResponse: true,
		SendPayload:   false,
	}

	fmt.Println("sc.Text?", sc.Text)

	// TODO refact c and commands variable names

	c := strings.Fields(sc.Text)
	var flags []string
	var noFlags []string

	for _, value := range c {
		if strings.HasPrefix(value, "-") {
			flags = append(flags, value)
		} else {
			noFlags = append(noFlags, value)
		}
	}

	if len(flags) > 0 {
		for _, flag := range flags {
			switch flag {
			case "-c":
				fmt.Println("send to payload channel")
				cp.Channel = fmt.Sprintf("#%v", sc.ChannelName)
				cp.SendPayload = true
				cp.SlashResponse = false
			case "-p":
				fmt.Println("send payload to user")
				cp.SendPayload = true
				cp.SlashResponse = false
			}
		}
	}

	c = noFlags

	fmt.Println("c strings?", c)

	// replace all underscores "_" with spaces " " in commands
	commands := make([]string, len(c))
	copy(commands, c)

	for i := 0; i < len(c); i++ {
		c[i] = strings.Replace(c[i], "_", " ", -1)
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
	if err != nil {
		return nil, err
	}

	// found boards return if only sent one command
	var boards []board
	json.Unmarshal(body, &boards)

	var responseString string

	// if the command is blank return possible commands(boards)
	if len(c) == 0 {

		// iterate over boards and create string to send
		for _, board := range boards {
			responseString += fmt.Sprint(board)
		}

		cp.Text = formatForSlack(c, responseString)

		return cp, nil
	}

	// if there is a command to access a board

	// first check if the command is a valid board
	var foundBoard board

	// check to see if the command passed is a list in the board
	for _, board := range boards {
		if strings.EqualFold(board.Name, c[0]) == true {
			foundBoard = board
			break
		}
	}

	// if nothing matches return
	if foundBoard.Name == "" {
		cp.Text = "invalid board name"
		return cp, nil
	}

	// first make sure the command is a valid list in the wiki board
	url = fmt.Sprintf(
		"https://api.trello.com/1/boards/%v/lists/?fields=name,idBoard&key=%v&token=%v",
		foundBoard.Id,
		os.Getenv("TRELLO_KEY"),
		os.Getenv("TRELLO_TOKEN"),
	)

	// fmt.Println("url: ", url)

	res, err = http.Get(url)
	body, _ = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var lists []list
	json.Unmarshal(body, &lists)

	// if the second command is black return lists
	if len(c) == 1 {

		fmt.Println("check command lenght", c)

		// iterate over boards and create string to send
		for _, list := range lists {
			responseString += fmt.Sprint(list)
		}

		cp.Text = formatForSlack(commands, responseString)

		return cp, nil
	}

	// if there is a command to access a list
	var foundList list

	// check to see if the command passed is a list in the board
	for _, list := range lists {
		if strings.EqualFold(list.Name, c[1]) == true {
			foundList = list
			break
		}
	}

	// if nothing matches return
	if foundList.Name == "" {
		cp.Text = "invalid board name"
		return cp, nil
	}

	// now look up cards in list
	url = fmt.Sprintf(
		"https://api.trello.com/1/lists/%v/?fields=name&cards=open&card_fields=name&key=%v&token=%v",
		foundList.Id,
		os.Getenv("TRELLO_KEY"),
		os.Getenv("TRELLO_TOKEN"),
	)
	res, err = http.Get(url)
	body, _ = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(body, &foundList)

	// iterate over boards and create string to send
	for _, card := range foundList.Cards {
		// check if they are urls
		card.URL = IsURL(card.Name)
		responseString += fmt.Sprint(card)
	}

	cp.Text = formatForSlack(commands, responseString)

	return cp, nil

}

func formatForSlack(c []string, s string) string {
	path := strings.Join(c, " ")
	return fmt.Sprintf("/fg %v ```%v```", path, s)
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
