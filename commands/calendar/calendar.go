package calendar

import (
    "encoding/json"
    "errors"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "os"
    "time" 

    "github.com/jesselucas/slackcmd/slack"
    "golang.org/x/net/context"
    "golang.org/x/oauth2"
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

func formatForSlack(s string) string {
    return fmt.Sprintf("```\n%v```", s)
}

func getAccessToken(clientID string, clientSecret string, refreshToken string) oauth2.Token {
    formValues := url.Values{
        "client_id": {clientID}, 
        "client_secret": {clientSecret},
        "refresh_token": {refreshToken},
        "grant_type": {"refresh_token"},
    }
    resp, _ := http.PostForm("https://accounts.google.com/o/oauth2/token", formValues)
    defer resp.Body.Close()

    bodyData, _ := ioutil.ReadAll(resp.Body)
    // body := string(bodyData)

    var token oauth2.Token
    json.Unmarshal(bodyData, &token)
    return token
}

func (cmd Command) Request(sc *slack.SlashCommand) (*slack.CommandPayload, error) {
    
    // Read the appropriate environment variables
    slackAPIKey := os.Getenv("SLACK_KEY_CALENDAR")
    clientID := os.Getenv("GOOGLE_CALENDAR_CLIENT_ID")
    clientSecret := os.Getenv("GOOGLE_CALENDAR_CLIENT_SECRET")
    refreshToken := os.Getenv("GOOGLE_CALENDAR_REFRESH_TOKEN")
    calendarID := os.Getenv("SLACK_CALENDAR_ID")

    if clientID == "" || clientSecret == "" || refreshToken == "" {
        err := errors.New("Server Configuration Error")
        return nil, err
    }

    // Verify the request is coming from Slack
    if sc.Token != slackAPIKey {
        err := errors.New("Unauthorized Slack")
        return nil, err
    }

    // Create initial payload
    payload := &slack.CommandPayload{
        Channel:       fmt.Sprintf("#%v", sc.ChannelName),
        Username:      "Calendar Bot",
        Emoji:         ":calendar:",
        SendPayload:   false,
        SlashResponse: true,
    }

    // Create the oauth config
    config := oauth2.Config {
        ClientID: clientID,
        ClientSecret: clientSecret,
        Endpoint: oauth2.Endpoint {
            AuthURL: "https://accounts.google.com/o/oauth2/auth", 
            TokenURL: "https://accounts.google.com/o/oauth2/token",
        },
        RedirectURL: "http://localhost",
    }

    // Create a client using our config, context, and access token
    context := context.Background()
    token := getAccessToken(clientID, clientSecret, refreshToken)
    client := config.Client(context, &token)

    // Get a calendar service
    service, err := calendar.New(client)
    if err != nil {
        err := errors.New("Unable to retrieve calendar client")
        return nil, err
    }

    // Setup the parameters for our calendar request.
    // We want to request all events from now until the end of today
    now := easyTime{time.Now()}
    timeMin := now.Format(time.RFC3339)
    timeMax := now.endOfDay().Format(time.RFC3339)

    // We want to request this information for a specific calendar ID,
    // which is indicated by the CALENDAR_ID environment variable
    events, err := service.Events.List(calendarID).ShowDeleted(false).SingleEvents(true).TimeMin(timeMin).TimeMax(timeMax).MaxResults(50).OrderBy("startTime").Do()
    if err != nil {
        err := errors.New("Unable to retrieve calendar events.")
        return nil, err
    }

    // Loop through the events received, and append them to the payload text
    payloadText := "Conference Room Schedule: " + time.Now().Format("Jan 2, 2006") + "\n"
    format := "03:04PM"
    if len(events.Items) > 0 {
        for _, i := range events.Items {
            // If the DateTime is an empty string the Event is an all-day Event.
            // So only Date is available.
            var timeString string
            if i.Start.DateTime != "" {
                start, startErr := time.Parse(time.RFC3339, i.Start.DateTime)
                end, endErr := time.Parse(time.RFC3339, i.End.DateTime)
                if startErr == nil && endErr == nil {
                    timeString = start.Local().Format(format) + " to " + end.Local().Format(format)
                } else {
                    timeString = "--------------"
                }
            } else {
                timeString = "All Day       "
            }

            payloadText += fmt.Sprintf("• [%v] <%v|%v>\n", timeString, i.HtmlLink, i.Summary)
        }
    } else {
        payloadText = "No upcoming events found.\n"
    }

    payload.Text = formatForSlack(payloadText)
    return payload, nil
}