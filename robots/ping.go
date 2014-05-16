package robots

import (
	"fmt"
)

type PingBot struct {
}

func init() {
	RegisterRobot("ping", func() (robot Robot) { return new(PingBot) })
}

func (p PingBot) Run(command *SlashCommand) (slashCommandImmediateReturn string) {
	go p.DeferredAction(command)
	return ""
}

func (p PingBot) DeferredAction(command *SlashCommand) {
	response := new(IncomingWebhook)
	response.Channel = command.Channel_ID
	response.Username = "Ping Bot"
	response.Text = fmt.Sprintf("@%s Pong!", command.User_Name)
	response.Icon_Emoji = ":ghost:"
	response.Unfurl_Links = true
	response.Parse = "full"
	MakeIncomingWebhookCall(response)
}

func (p PingBot) Description() (description string) {
	return "Ping bot!\n\tUsage: /ping\n\tExpected Response: @user: Pong!"
}
