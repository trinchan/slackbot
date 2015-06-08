slackbot
===========
Simple, pluggable bot framework for [Slack](https://www.slack.com) chat.  

Installation
============
You can grab the source code using `go get github.com/trinchan/slackbot` and install like usual. `go install github.com/trinchan/slackbot`

Setup
=====
Slackbot uses environment variables for all configuration, with domain wide variables prefixed by the domain name in title case (e.g. `MYDOMAIN_IN_URL` for MyDomain's incoming webhook URL). This makes it easy to support multiple domains and deploy to Heroku. Make sure to set a `PORT` environment variable defining what port to run your bot on.

An example environment variable list:
```
BIJIN_TIMEZONE=Asia/Tokyo
MYDOMAIN_OUT_TOKEN=AVerySecretToken
MYDOMAIN_IN_URL=https://hooks.slack.com/services/AAAAAAAAA/BBBBBBBBB/AnotherVerySecretToken
PORT=5000
```

###Setup an Incoming Webhook
If you don't already have an [Incoming Webhook](https://my.slack.com/services/new/incoming-webhook) setup for your bot, you'll want to start here.  Set one up (make sure it's enabled) and don't be afraid to read through the setup instructions.  You're after the "Webhook URL" that slack generates for you.
```
https://hooks.slack.com/services/AAAAAAAAA/BBBBBBBBB/YourSecretToken123456789
```

For each domain, set an environment variable `DOMAIN_IN_URL` to this URL.

###Send messages to your bot
This framework can respond to "slash commands" and "outgoing webhooks"  If you want users to be able to silently type `/ping`, and have the ping-bot respond in their channel, then you'll want to set up "slash commands".  Each bot will need it's own command setup.  The other option is to configure an outgoing webhook with a symbol for the prefix. Exe: `!ping`.  This option only requires one configuration, but the commands will be entered into the channel as regular messages.

#####Configuring an Outgoing Webhook
I use an [Outgoing Webhook](https://my.slack.com/services/new/outgoing-webhook)

1. Add a new Outgoing Webhook Integration.
2. Here, you can tell slack to ONLY pay attention to a specific channel, or to simply listen to all public channels.  Outgoing Webhooks can not listen to private channels/direct messages.
3. For each domain, set an environment variable `DOMAIN_OUT_TOKEN` to your integration's token. This is used to verify payloads come from Slack.
4. The {trigger_word} should only be one character (preferrably a symbol, such as ! or ?) and typing `{trigger_word}ping` will trigger the Ping bot.  
5. The URL should follow the following format: `your_address.com:port/slack_hook` (no trailing /)  
The bot will respond to commands of the form `{trigger_word}bot param param param` in the specified channels

#####Configuring Slash Commands
Alternatively, each bot you make can respond to a corresponding [Slash Command](https://my.slack.com/services/new/slash-commands).

1. Add a new slash command, use the [bot's name](https://github.com/trinchan/slackbot/tree/master/robots) as the name of the command.  
2. The URL should follow the following format: `your_address.com:port/slack` (no trailing /)  
3. You want to use POST.  
4. For each bot, set an environment variable `BOTNAME_SLACK_TOKEN` to your slash command's token. This is used to verify payloads come from Slack.
5. Repeat for each bot you want to enable.

The bot will respond to commands of the form `/bot param param param`

###Configuring Heroku
After setting up the proper environment variables, deploying to heroku should be as simple using the [heroku-go-buildpack](https://github.com/trinchan/heroku-buildpack-go) with a one line modification to run `go generate ./...` before installing to generate the plugin import file.

Adding Bots
===========
1. Create a new package and implement the [Robot](https://github.com/trinchan/slackbot/tree/master/robots/robot.go) interface.
2. In [importer/importer.sh](https://github.com/trinchan/slackbot/tree/master/importer/importer.sh), add the path to your package to the robots array.
3. Run `go generate ./...` to generate `init.go` which will import your bot.
4. Rebuild slackbot and deploy.

If you use [Sublime Text](http://www.sublimetext.com/) or another editor which supports snippets for development, then you can simply add and use the included snippet to easily generate a template Robot based on the filename. Otherwise, refer to the template below.

```go
package test

import "github.com/trinchan/slackbot/robots"

type bot struct{}

// Registers the bot with the server for command /test.
func init() {
	r := &bot{}
	robots.RegisterRobot("test", r)
}

// All Robots must implement a Run command to be executed when the registered command is received.
func (r bot) Run(p *robots.Payload) string {
	// If you (optionally) want to do some asynchronous work (like sending API calls to slack)
	// you can put it in a go routine like this
	go r.DeferredAction(p)
	// The string returned here will be shown only to the user who executed the command
	// and will show up as a message from slackbot.
	return "Text to be returned only to the user who made the command."
}

func (r bot) DeferredAction(p *robots.Payload) {
	// Let's use the IncomingWebhook struct defined in payload.go to form and send an
	// IncomingWebhook message to slack that can be seen by everyone in the room. You can
	// read the Slack API Docs (https://api.slack.com/) to know which fields are required, etc.
	// You can also see what data is available from the command structure in definitions.go
	response := &robots.IncomingWebhook{
		Channel:     p.ChannelID,
		Username:    "Test Bot",
		Text:        "Hi there!",
		IconEmoji:   ":ghost:",
		UnfurlLinks: true,
		Parse:       robots.ParseStyleFull,
	}
	response.Send()
}

func (r bot) Description() string {
	// In addition to a Run method, each Robot must implement a Description method which
	// is just a simple string describing what the Robot does. This is used in the included
	// /c command which gives users a list of commands and descriptions
	return "This is a description for TestBot which will be displayed on /help"
}

```

If you are using [Slash Commands](https://my.slack.com/services/new/slash-commands), you'll need to add a new slash command integration for each bot you add.

Running
=======
If you are not using Heroku, making a `env.sh` file which exports your environment variables and running slackbot via
```
go generate ./... && source env.sh && slackbot
```
makes for a convenient one-liner.

If you see output similar to below and you have the commands enabled in your Slack integration, you're ready to go!
```
2015/06/07 23:26:56 Registered: wiki
2015/06/07 23:26:56 Registered: store
2015/06/07 23:26:56 Registered: roll
2015/06/07 23:26:56 Registered: ping
2015/06/07 23:26:56 Registered: nihongo
2015/06/07 23:26:56 Registered: help
2015/06/07 23:26:56 Registered: decide
2015/06/07 23:26:56 Registered: bot
2015/06/07 23:26:56 Registered: bijin
2015/06/07 23:26:56 Starting HTTP server on 13748
```
