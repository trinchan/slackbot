package robots

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type WikiBot struct {
}

func init() {
	RegisterRobot("/wiki", func() (robot Robot) { return new(WikiBot) })
}

func (w WikiBot) Run(command *SlashCommand) (slashCommandImmediateReturn string) {
	go w.DeferredAction(command)
	return ""
}

func (w WikiBot) DeferredAction(command *SlashCommand) {
	text := strings.TrimSpace(command.Text)
	if text != "" {
		response := new(IncomingWebhook)
		response.Channel = command.Channel_ID
		response.Username = "Wiki Bot"
		response.Text = fmt.Sprintf("@%s: Searching google for wikipedias related to: %s", command.User_Name, text)
		response.Icon_Emoji = ":ghost:"
		response.Unfurl_Links = true
		response.Parse = "full"
		MakeIncomingWebhookCall(response)
		resp, err := http.Get(fmt.Sprintf("http://www.google.com/search?q=site:*.wikipedia.org+%s&btnI", text))
		if err != nil {
			response.Text = fmt.Sprintf("@%s: %s", command.User_Name, "Error getting wikipedia link from google :(")
		} else {
			response.Text = fmt.Sprintf("@%s: %s", command.User_Name, resp.Request.URL.String())
		}
		MakeIncomingWebhookCall(response)
	}
}

func (w WikiBot) Description() (description string) {
	return "Wiki bot!\n\tUsage: /wiki\n\tExpected Response: @user: Link to wikipedia article!"
}
