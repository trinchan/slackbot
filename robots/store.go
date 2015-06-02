package robots

import "fmt"

type StoreBot struct {
}

type StoreConfiguration struct {
}

var StoreConfig = new(StoreConfiguration)

// Loads the config file and registers the bot with the server for command /store.
func init() {
	s := &StoreBot{}
	RegisterRobot("store", s)
}

// All Robots must implement a Run command to be executed when the registered command is received.
func (r StoreBot) Run(p *Payload) (slashCommandImmediateReturn string) {
	// If you (optionally) want to do some asynchronous work (like sending API calls to slack)
	// you can put it in a go routine like this
	go r.DeferredAction(p)
	// The string returned here will be shown only to the user who executed the command
	// and will show up as a message from slackbot.
	return ""
}

func (r StoreBot) DeferredAction(p *Payload) {
	// Let's use the IncomingWebhook struct defined in definitions.go to form and send an
	// IncomingWebhook message to slack that can be seen by everyone in the room. You can
	// read the Slack API Docs (https://api.slack.com/) to know which fields are required, etc.
	// You can also see what data is available from the command structure in definitions.go
	response := &IncomingWebhook{
		Domain:      p.TeamDomain,
		Channel:     p.ChannelID,
		Username:    "Store Bot",
		IconEmoji:   ":famima:",
		Text:        fmt.Sprintf("@group @%s wants to go to the store", p.UserName),
		UnfurlLinks: true,
		Parse:       ParseStyleFull,
	}
	response.Send()
}

func (r StoreBot) Description() (description string) {
	// In addition to a Run method, each Robot must implement a Description method which
	// is just a simple string describing what the Robot does. This is used in the included
	// /c command which gives users a list of commands and descriptions
	return "This is a description for StoreBot which will be displayed on /c"
}
