slackbot
===========

A framework for building slash commands for Slack in Go

Dependencies
============
Schema  - `go get github.com/gorilla/schema`

Installation
============
You can grab the source code using `go get github.com/trinchan/slackbot` and install like usual.

Setup
=====
Create a config file (config.json) to the executing directory with the following format:

```json
{
    "domain": "{YOUR_SLACK_DOMAIN}",
    "port": {PORT_FOR_BOT},
    "token": "{YOUR_SLACK_INCOMING_WEBHOOK_TOKEN}"
}
```

Make sure you have [Incoming Webhooks](https://my.slack.com/services/new/incoming-webhook) enabled and you are using that integration token for your config.

For each slash command (including the default commands!), be sure to add a corresponding entry in [Slash Commands](https://my.slack.com/services/new/slash-commands) to POST to server:port/slack of your slackbot setup. Note no trailing slash after /slack.

Adding Bots
===========
Create a new file in the robots directory and follow the template below
```go
package robots

type ExampleBot struct {
}
// Registers the bot with the server for command /example.
func init() {
	RegisterRobot("/example", func() (robot Robot) { return new(ExampleBot) })
}

// All Robots must implement a Run command to be executed when the registered command is received.
func (e ExampleBot) Run(command *SlashCommand) (slashCommandImmediateReturn string) {
	// If you (optionally) want to do some asynchronous work (like sending API calls to slack)
	// you can put it in a go routine like this 
	go e.DeferredAction(command)
	// The string returned here will be shown only to the user who executed the command
	// and will show up as a message from slackbot.
	return "Text to be returned only to the user who made the command."
}

func (e ExampleBot) DeferredAction(command *SlashCommand) {
    // Let's use the IncomingWebhook struct defined in definitions.go to form and send an 
    // IncomingWebhook message to slack that can be seen by everyone in the room. You can 
    // read the Slack API Docs (https://api.slack.com/) to know which fields are required, etc. 
    // You can also see what data is available from the command structure in definitions.go
	response := new(IncomingWebhook)
	response.Channel = command.Channel_ID
	response.Username = "Example Bot"
	response.Text = "Hi there!"
	response.Icon_Emoji = ":ghost:"
	response.Unfurl_Links = true
	response.Parse = "full"
	MakeIncomingWebhookCall(response)
}

func (e ExampleBot) Description() (description string) {
	// In addition to a Run method, each Robot must implement a Description method which
	// is just a simple string describing what the Robot does. This is used in the included
	// /c command which gives users a list of commands and descriptions
	return "This is a description for example bot which will"
}
```

Now just add a corresponding entry in [Slash Commands](https://my.slack.com/services/new/slash-commands) to POST to server:port/slack of your slackbot setup. Note no trailing slash after /slack.

Running
=======
If you see output similar to below and you have the commands enabled in your Slack integration, you're ready to go!
```
2014/02/18 10:55:07 Registered: /decide
2014/02/18 10:55:07 Registered: /ping
2014/02/18 10:55:07 Registered: /c
2014/02/18 10:55:07 Registered: /roll
2014/02/18 10:55:07 Starting HTTP server on 8888
```
