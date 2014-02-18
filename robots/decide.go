package robots

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type DecideBot struct {
}

func init() {
	RegisterRobot("/decide", func() (robot Robot) { return new (DecideBot) })
}

func (d DecideBot) Run(command *SlashCommand) (slashCommandImmediateReturn string) {
	go d.DeferredAction(command)
	text := strings.TrimSpace(command.Text)
	if text == "" {
		return "I need something to decide on!"
	} else {
		return ""
	}
}

func (d DecideBot) DeferredAction(command *SlashCommand) {
	response := new (IncomingWebhook)
	response.Channel = command.Channel_ID
	response.Username = "Fate Bot"
	response.Icon_Emoji   = ":ghost:"
	response.Unfurl_Links = true
	response.Parse = "full"
	text := strings.TrimSpace(command.Text)
	if text != "" {
		split := strings.Split(text, " ")
		response.Text = fmt.Sprintf("@%s: Deciding between: (%s)", command.User_Name, strings.Join(split, ", "))
		MakeIncomingWebhookCall(response)
		response.Text = fmt.Sprintf("@%s: Decided on: %s", command.User_Name, Decide(split))
		MakeIncomingWebhookCall(response)
	}
}

func (d DecideBot) Description() (description string) {
	return "Decides your fate!\n\tUsage: /decide Life Death ...\n\tExpected Response: Deciding on (Life, Death, ...)\n\tDecided on Life!"
}

func Decide(Fates []string) (result string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := len(Fates)
	if n > 0 {
		return Fates[r.Intn(n)]
	} else {
		return fmt.Sprintf("Error")
	}
}