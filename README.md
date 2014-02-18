slackbot-go
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
