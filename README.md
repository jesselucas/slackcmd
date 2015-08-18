## SlackCMD
Go app to create custom *Slack* slash commands and easily hook into bots webhook.

## Commands Package
Every package must use it's own slash command token. Ex. `SLACK_KEY_COMMAND`

We use this to verify the request came from your Slack team. -  - https://api.slack.com/slash-commands 

### Trello
Slack token: `SLACK_KEY_TRELLO`

Reads boards from a user's organization on Trello. To set this up for your account you must first get a a key and token from Trello - https://trello.com/docs/gettingstarted/

These are read in through environment variables `os.Getenv("TRELLO_KEY")` and `os.Getenv("TRELLO_TOKEN")`. 

`trelloOrg` sets the organization name you want to access

### Beats1
Slack token: `SLACK_KEY_BEATS1`

Accesses the `beats1plays` twitter stream and reads the latest tweet. Which is the current song being played on Beats1.

It uses twitter's api: https://dev.twitter.com/rest/public

```
consumerKey := os.Getenv("TWITTER_CONSUMER_KEY")
consumerSecret := os.Getenv("TWITTER_CONSUMER_SECRET")
accessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
accessTokenSecret := os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")
```

### Calendar
Slack token: SLACK_KEY_CALENDAR