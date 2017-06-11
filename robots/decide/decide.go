package decide

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/trinchan/slackbot/robots"
)

type bot struct {
	Random *rand.Rand
}

func init() {
	d := &bot{}
	d.Random = rand.New(rand.NewSource(time.Now().UnixNano()))
	robots.RegisterRobot("decide", d)
}

func (d bot) Run(p *robots.Payload) (slashCommandImmediateReturn string) {
	go d.DeferredAction(p)
	text := strings.TrimSpace(p.Text)
	if text == "" {
		return "I need something to decide on!"
	}
	return ""
}

func (d bot) DeferredAction(p *robots.Payload) {
	response := robots.IncomingWebhook{
		Domain:      p.TeamDomain,
		Channel:     p.ChannelID,
		Username:    "Fate Bot",
		IconEmoji:   ":ghost:",
		UnfurlLinks: true,
		Parse:       robots.ParseStyleFull,
	}
	text := strings.TrimSpace(p.Text)
	if text != "" {
		split := strings.Split(text, ",")
		response.Text = fmt.Sprintf("@%s: Deciding between: (%s) -> %s", p.UserName, strings.Join(split, ","), d.decide(split))
		response.Send()
	}
}

func (d bot) Description() (description string) {
	return "Decides your fate!\n\tUsage: /decide Life Death ...\n\tExpected Response: Deciding on (Life, Death, ...)\n\tDecided on Life!"
}

func (d bot) decide(fates []string) (result string) {
	n := len(fates)
	if n > 0 {
		return fates[d.Random.Intn(n)]
	}
	return fmt.Sprintf("Error")
}
