slackbot
===========

A framework for building slash commands for Slack in Go

Add a config file (config.json) to the root directory of the project with the following format:

```json
{
    "domain": "{YOUR_SLACK_DOMAIN}",
    "port": {PORT_FOR_BOT},
    "token": "{YOUR_SLACK_INCOMING_WEBHOOK_TOKEN}"
}
```

Dependencies
============
Schema  - `go get github.com/gorilla/schema`

Installation
============
`go get github.com/trinchan/slackbot`

Make sure you have [Incoming Webhooks](https://my.slack.com/services/new/incoming-webhook) enabled and you are using that integration token for your config.

For each slash command, be sure to add a corresponding entry in [Slash Commands](https://my.slack.com/services/new/slash-commands) to POST to server:port/slack of your slackbot setup. Note no trailing slash after /slack. 
