package robots

import (
	"fmt"

	"github.com/trinchan/slackbot/robots"
)

type bot struct{}

// Loads the config file and registers the bot with the server for command /store.
func init() {
	s := &bot{}
	robots.RegisterRobot("store", s)
}

// All Robots must implement a Run command to be executed when the registered command is received.
func (r bot) Run(p *robots.Payload) (slashCommandImmediateReturn string) {
	// If you (optionally) want to do some asynchronous work (like sending API calls to slack)
	// you can put it in a go routine like this
	go r.DeferredAction(p)
	// The string returned here will be shown only to the user who executed the command
	// and will show up as a message from slackbot.
	return ""
}

func (r bot) DeferredAction(p *robots.Payload) {
	// Let's use the IncomingWebhook struct defined in definitions.go to form and send an
	// IncomingWebhook message to slack that can be seen by everyone in the room. You can
	// read the Slack API Docs (https://api.slack.com/) to know which fields are required, etc.
	// You can also see what data is available from the command structure in definitions.go
	response := &robots.IncomingWebhook{
		Domain:      p.TeamDomain,
		Channel:     p.ChannelID,
		Username:    "Store Bot",
		IconEmoji:   ":famima:",
		Text:        fmt.Sprintf("@group @%s wants to go to the store", p.UserName),
		UnfurlLinks: true,
		Parse:       robots.ParseStyleFull,
	}
	response.Send()
}

func (r bot) Description() (description string) {
	// In addition to a Run method, each Robot must implement a Description method which
	// is just a simple string describing what the Robot does. This is used in the included
	// /c command which gives users a list of commands and descriptions
	return "This is a description for Bot which will be displayed on /c"
}
