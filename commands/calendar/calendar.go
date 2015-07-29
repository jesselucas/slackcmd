package calendar

import (
    "encoding/json"
    "errors"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "os/user"
    "os"
    "path/filepath"
    "time" 

    "github.com/jesselucas/slackcmd/slack"
    "golang.org/x/net/context"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "google.golang.org/api/calendar/v3"
)

type easyTime struct {
    time.Time
}

func (t *easyTime) beginningOfHour() time.Time {
    return t.Truncate(time.Hour)
}

func (t *easyTime) beginningOfDay() time.Time {
    d := time.Duration(-t.Hour()) * time.Hour
    return t.beginningOfHour().Add(d)
}

func (t *easyTime) endOfDay() time.Time {
    return t.beginningOfDay().Add(24*time.Hour - time.Nanosecond)
}

type Command struct {
}

func (cmd Command) Request(sc *slack.SlashCommand) (*slack.CommandPayload, error) {
    // Verify the request is coming from Slack
    slackAPIKey := os.Getenv("SLACK_KEY_CALENDAR")
    if sc.Token != slackAPIKey {
        err := errors.New("Unauthorized Slack")
        return nil, err
    }

    // create payload
    payload := &slack.CommandPayload{
        Channel:       fmt.Sprintf("#%v", sc.ChannelName),
        Username:      "Calendar Bot",
        Emoji:         ":calendar:",
        SendPayload:   false,
        SlashResponse: true,
    }

    context := context.Background()

    //Get the client_secret json file path
    clientSecretPath := os.Getenv("CALENDAR_CLIENT_SECRET_PATH")
    jsonConfig, err := ioutil.ReadFile(clientSecretPath)
    if err != nil {
        err := errors.New("Unable to read client secret file")
        return nil, err
    }

    //Get the calendar configuration from the json data
    config, err := google.ConfigFromJSON(jsonConfig, calendar.CalendarReadonlyScope)
    if err != nil {
        err := errors.New("Unable to parse client secret file to config")
        return nil, err
    }

    //Get a calendar client based on our configuration
    client := getClient(context, config)
    service, err := calendar.New(client)
    if err != nil {
        err := errors.New("Unable to retrieve calendar client")
        return nil, err
    }

    //Setup the parameters for our calendar request.
    //We want to request all events from now until the end of today
    now := easyTime{time.Now()}
    timeMin := now.Format(time.RFC3339)
    timeMax := now.endOfDay().Format(time.RFC3339)

    //We want to request this information for a specific calendar ID,
    //which is indicated by the CALENDAR_ID environment variable
    calendarID := os.Getenv("SLACK_CALENDAR_ID")
    events, err := service.Events.List(calendarID).ShowDeleted(false).SingleEvents(true).TimeMin(timeMin).TimeMax(timeMax).MaxResults(50).OrderBy("startTime").Do()
    if err != nil {
        err := errors.New("Unable to retrieve calendar events.")
        return nil, err
    }

    //Loop through the events received, and append them to the payload text
    payloadText := ""
    if len(events.Items) > 0 {
        for _, i := range events.Items {
            // If the DateTime is an empty string the Event is an all-day Event.
            // So only Date is available.
            var timeString string
            if i.Start.DateTime != "" {
                timeString = i.Start.DateTime + " to " + i.End.DateTime
            } else {
                timeString = "All Day"
            }

            payloadText += fmt.Sprintf("%s (%s) %v\n", i.Summary, timeString, i.HtmlLink)
        }
    } else {
        payloadText = "No upcoming events found.\n"
    }

    fmt.Println(payloadText)
    payload.Text = formatForSlack(payloadText)
    return payload, nil

}

func formatForSlack(s string) string {
    return fmt.Sprintf("```\n%v```", s)
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(context context.Context, config *oauth2.Config) *http.Client {
    cacheFile, err := tokenCacheFile()
    if err != nil {
        log.Fatalf("Unable to get path to cached credential file.  %v", err)        
    }

    token, err := tokenFromFile(cacheFile)
    if err != nil {
        token = getTokenFromWeb(config)
        saveToken(cacheFile, token)
    }    
    
    return config.Client(context, token)
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
    user, err := user.Current()
    if err != nil {
        return "", err
    }

    tokenCacheDir := filepath.Join(user.HomeDir, ".credentials")
    os.MkdirAll(tokenCacheDir, 0700)
    return filepath.Join(tokenCacheDir, url.QueryEscape("calendar-api-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
    f, err := os.Open(file)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    token := &oauth2.Token{}
    err  = json.NewDecoder(f).Decode(token)
    return token, err
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
    authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
    fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)
    
    var code string
    if _, err := fmt.Scan(&code); err != nil {
        log.Fatalf("Unable to read authorization code %v", err)
    }    

    token, err := config.Exchange(oauth2.NoContext, code)
    if err != nil {
        log.Fatalf("Unable to retrieve token from web %v", err)
    }
    return token
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
    fmt.Printf("Saving credential file to %s\n", file)
    f, err := os.Create(file)
    if err != nil {
        log.Fatalf("Unable to cache oauth token: %v", err)
    }
    defer f.Close()

    json.NewEncoder(f).Encode(token)
}