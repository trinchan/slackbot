package robots

import (
	"fmt"
)

type PingBot struct {
}

func init() {
	p := &PingBot{}
	RegisterRobot("ping", p)
}

func (pb PingBot) Run(p *Payload) (slashCommandImmediateReturn string) {
	go pb.DeferredAction(p)
	return ""
}

func (pb PingBot) DeferredAction(p *Payload) {
	response := &IncomingWebhook{
		Domain:      p.TeamDomain,
		Channel:     p.ChannelID,
		Username:    "Ping Bot",
		Text:        fmt.Sprintf("@%s Pong!", p.UserName),
		IconEmoji:   ":ghost:",
		UnfurlLinks: true,
		Parse:       ParseStyleFull,
	}
	response.Send()
}

func (pb PingBot) Description() (description string) {
	return "Ping bot!\n\tUsage: /ping\n\tExpected Response: @user: Pong!"
}
