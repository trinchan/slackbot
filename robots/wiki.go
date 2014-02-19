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
		resp, err = http.Get(fmt.Sprintf("http://www.google.com/search?q=wikipedia+%s&btnI", strings.Replace(text, " ", "+", -1)))
		if err != nil {
			log.Printf("Error getting wiki link from google!")
			return
		}
		response := new(IncomingWebhook)
		response.Channel = command.Channel_ID
		response.Username = "Wiki Bot"
		response.Text = fmt.Sprintf("@%s: %s!", command.User_Name, resp.Request.URL.String())
		response.Icon_Emoji = ":ghost:"
		response.Unfurl_Links = true
		response.Parse = "full"
		MakeIncomingWebhookCall(response)
	}
}

func (w WikiBot) Description() (description string) {
	return "Wiki bot!\n\tUsage: /wiki\n\tExpected Response: @user: Link to wikipedia article!"
}
