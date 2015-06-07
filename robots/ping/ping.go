package robots

import (
	"fmt"

	"github.com/trinchan/slackbot/robots"
)

type bot struct{}

func init() {
	p := &bot{}
	robots.RegisterRobot("ping", p)
}

func (pb bot) Run(p *robots.Payload) (slashCommandImmediateReturn string) {
	go pb.DeferredAction(p)
	return ""
}

func (pb bot) DeferredAction(p *robots.Payload) {
	response := &robots.IncomingWebhook{
		Domain:      p.TeamDomain,
		Channel:     p.ChannelID,
		Username:    "Ping Bot",
		Text:        fmt.Sprintf("@%s Pong!", p.UserName),
		IconEmoji:   ":ghost:",
		UnfurlLinks: true,
		Parse:       robots.ParseStyleFull,
	}
	response.Send()
}

func (pb bot) Description() (description string) {
	return "Ping bot!\n\tUsage: /ping\n\tExpected Response: @user: Pong!"
}
