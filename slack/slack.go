package slack

import (
	"strings"
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

// CommandPayload https://api.slack.com/docs/formatting
type CommandPayload struct {
	Channel       string       `json:"channel"`
	Username      string       `json:"username"`
	Emoji         string       `json:"icon_emoji"`
	EmojiURL      string       `json:"icon_url"`
	Text          string       `json:"text"`
	Attachments   []Attachment `json:"attachments"`
	UnfurlMedia   bool         `json:"unfurl_media"`
	UnfurlLinks   bool         `json:"unfurl_links"`
	Parse         string       `json:"parse"`
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

// Takes Slack slash command text and parses out any flags
// Ex. "golang links -c" returns "golang links"
func SeparateFlags(t string) (c string, f []string) {
	var splitFlags = strings.Fields(t)
	var parsedFlags []string
	var trimmedFlag string
	var parsedCommands []string

	// first seperate the string into an slice of strings
	for _, value := range splitFlags {
		// Test for flags and remove prefix.
		// Check -- first since - will always find --
		if strings.HasPrefix(value, "-") || strings.HasPrefix(value, "—") {
			if strings.HasPrefix(value, "——") {
				trimmedFlag = strings.TrimPrefix(value, "——")
				parsedFlags = append(parsedFlags, trimmedFlag)
			} else if strings.HasPrefix(value, "--") {
				trimmedFlag = strings.TrimPrefix(value, "--")
				parsedFlags = append(parsedFlags, trimmedFlag)
			} else if strings.HasPrefix(value, "—") {
				trimmedFlag = strings.TrimPrefix(value, "—")
				parsedFlags = append(parsedFlags, trimmedFlag)
			} else {
				trimmedFlag = strings.TrimPrefix(value, "-")
				parsedFlags = append(parsedFlags, trimmedFlag)
			}

		} else {
			parsedCommands = append(parsedCommands, value)
		}
	}

	// Return string without flags and flags
	return strings.Join(parsedCommands, " "), parsedFlags
}

func SanitizeString(s string) string {
	// 	& replaced with &amp;
	// < replaced with &lt;
	// > replaced with &gt;

	strings.Replace(s, "&", "&amp;", -1)
	strings.Replace(s, "<", "&lt;", -1)
	strings.Replace(s, ">", "&gt;/", -1)

	return s
}
