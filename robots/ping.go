package robots

import (
	"fmt"
)

type PingBot struct {
}

func init() {
	RegisterRobot("ping", new(PingBot))
}

func (pb PingBot) Run(p *Payload) (slashCommandImmediateReturn string) {
	go pb.DeferredAction(p)
	return ""
}

func (pb PingBot) DeferredAction(p *Payload) {
	response := &IncomingWebhook{
		Channel:     p.ChannelID,
		Username:    "Ping Bot",
		Text:        fmt.Sprintf("@%s Pong!", p.UserName),
		IconEmoji:   ":ghost:",
		UnfurlLinks: true,
		Parse:       ParseStyleFull,
	}
	MakeIncomingWebhookCall(response)
}

func (pb PingBot) Description() (description string) {
	return "Ping bot!\n\tUsage: /ping\n\tExpected Response: @user: Pong!"
}
