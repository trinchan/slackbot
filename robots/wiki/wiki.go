package robots

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/trinchan/slackbot/robots"
)

type bot struct{}

func init() {
	w := &bot{}
	robots.RegisterRobot("wiki", w)
}

func (w bot) Run(p *robots.Payload) (slashCommandImmediateReturn string) {
	go w.DeferredAction(p)
	return ""
}

func (w bot) DeferredAction(p *robots.Payload) {
	text := strings.TrimSpace(p.Text)
	if text != "" {
		response := &robots.IncomingWebhook{
			Domain:      p.TeamDomain,
			Channel:     p.ChannelID,
			Username:    "Wiki Bot",
			Text:        fmt.Sprintf("@%s: Searching google for wikis relating to: %s", p.UserName, text),
			IconEmoji:   ":ghost:",
			UnfurlLinks: true,
			Parse:       robots.ParseStyleFull,
		}

		go response.Send()
		resp, err := http.Get(fmt.Sprintf("http://www.google.com/search?q=(site:en.wikipedia.org+OR+site:ja.wikipedia.org)+%s&btnI", url.QueryEscape(text)))
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			message := fmt.Sprintf("ERROR: Non-200 Response from Google: %s", resp.Status)
			log.Println(message)
			response.Text = fmt.Sprintf("@%s: %s", p.UserName, message)
		} else if err != nil {
			response.Text = fmt.Sprintf("@%s: %s", p.UserName, "Error getting wikipedia link from google :(")
		} else {
			response.Text = fmt.Sprintf("@%s: %s", p.UserName, resp.Request.URL.String())
		}
		response.Send()
	}
}

func (w bot) Description() (description string) {
	return "Wiki bot!\n\tUsage: /wiki\n\tExpected Response: @user: Link to wikipedia article!"
}
