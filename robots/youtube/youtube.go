package robots

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/trinchan/slackbot/robots"
)

type bot struct{}

func init() {
	y := &bot{}
	robots.RegisterRobot("youtube", y)
}

// All Robots must implement a Run command to be executed when the registered command is received.
func (r bot) Run(p *robots.Payload) string {
	// If you (optionally) want to do some asynchronous work (like sending API calls to slack)
	// you can put it in a go routine like this
	go r.DeferredAction(p)
	// The string returned here will be shown only to the user who executed the command
	// and will show up as a message from slackbot.
	return ""
}

type youTubeVideoFeedResults struct {
	Feed youTubeFeed `json:"feed"`
}

type youTubeFeed struct {
	Entries []youTubeEntry `json:"entry,omitempty"`
}

type youTubeEntry struct {
	Title   youTubeTitle   `json:"title"`
	Content youTubeContent `json:"content"`
	Link    []youTubeLink  `json:"link"`
}

type youTubeTitle struct {
	String string `json:"$t"`
}

type youTubeContent struct {
	String string `json:"$t"`
}

type youTubeLink struct {
	Relative string `json:"rel"`
	Type     string `json:"type"`
	Href     string `json:"href"`
}

func (r bot) DeferredAction(p *robots.Payload) {
	text := strings.TrimSpace(p.Text)
	if text != "" {
		response := &robots.IncomingWebhook{
			Domain:      p.TeamDomain,
			Channel:     p.ChannelID,
			Username:    "YouTube Bot",
			Text:        fmt.Sprintf("@%s: Searching youtube for %s", p.UserName, text),
			IconEmoji:   ":ghost:",
			UnfurlLinks: true,
		}

		go response.Send()
		resp, err := http.Get(fmt.Sprintf("https://gdata.youtube.com/feeds/api/videos?q=%s&orderBy=relevance&alt=json&max-results=1", url.QueryEscape(text)))
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			message := fmt.Sprintf("ERROR: Non-200 Response from YouTube: %s", resp.Status)
			log.Println(message)
			response.Text = fmt.Sprintf("@%s: %s", p.UserName, message)
		} else if err != nil {
			response.Text = fmt.Sprintf("@%s: %s", p.UserName, "Error getting YouTube video :(")
		} else {
			results := youTubeVideoFeedResults{}
			r, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				response.Text = fmt.Sprintf("@%s: %s", p.UserName, "Error getting YouTube video :(")
			} else {
				err = json.Unmarshal(r, &results)
				if err != nil {
					response.Text = fmt.Sprintf("@%s: %s", p.UserName, "Error getting YouTube video :(")
				}
			}
			if len(results.Feed.Entries) > 0 {
				response.Text = fmt.Sprintf("<@%s|%s>  <%s|%s - %s> ", p.UserID, p.UserName, results.Feed.Entries[0].Link[0].Href, results.Feed.Entries[0].Title.String, results.Feed.Entries[0].Content.String)
			} else {
				response.Text = fmt.Sprintf("@%s: %s", p.UserName, "No YouTube videos for that search :(")
			}
		}
		response.Send()
	}
}

func (r bot) Description() (description string) {
	// In addition to a Run method, each Robot must implement a Description method which
	// is just a simple string describing what the Robot does. This is used in the included
	// /c command which gives users a list of commands and descriptions
	return "Gets the most relevant YouTube result"
}
