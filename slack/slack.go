package slack

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
