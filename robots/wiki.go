package robots

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type WikiBot struct {
}

func init() {
	w := &WikiBot{}
	RegisterRobot("wiki", w)
}

func (w WikiBot) Run(p *Payload) (slashCommandImmediateReturn string) {
	go w.DeferredAction(p)
	return ""
}

func (w WikiBot) DeferredAction(p *Payload) {
	text := strings.TrimSpace(p.Text)
	if text != "" {
		response := &IncomingWebhook{
			Domain:      p.TeamDomain,
			Channel:     p.ChannelID,
			Username:    "Wiki Bot",
			Text:        fmt.Sprintf("@%s: Searching google for wikis relating to: %s", p.UserName, text),
			IconEmoji:   ":ghost:",
			UnfurlLinks: true,
			Parse:       ParseStyleFull,
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

func (w WikiBot) Description() (description string) {
	return "Wiki bot!\n\tUsage: /wiki\n\tExpected Response: @user: Link to wikipedia article!"
}
