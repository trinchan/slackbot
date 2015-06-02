package robots

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type YoutubeBot struct {
}

type YoutubeConfiguration struct {
}

var YoutubeConfig = new(YoutubeConfiguration)

func init() {
	y := &YoutubeBot{}
	RegisterRobot("youtube", y)
}

// All Robots must implement a Run command to be executed when the registered command is received.
func (r YoutubeBot) Run(p *Payload) string {
	// If you (optionally) want to do some asynchronous work (like sending API calls to slack)
	// you can put it in a go routine like this
	go r.DeferredAction(p)
	// The string returned here will be shown only to the user who executed the command
	// and will show up as a message from slackbot.
	return ""
}

type YouTubeVideoFeedResults struct {
	Feed YouTubeFeed `json:"feed"`
}

type YouTubeFeed struct {
	Entries []YouTubeEntry `json:"entry,omitempty"`
}

type YouTubeEntry struct {
	Title   YouTubeTitle   `json:"title"`
	Content YouTubeContent `json:"content"`
	Link    []YouTubeLink  `json:"link"`
}

type YouTubeTitle struct {
	String string `json:"$t"`
}

type YouTubeContent struct {
	String string `json:"$t"`
}

type YouTubeLink struct {
	Relative string `json:"rel"`
	Type     string `json:"type"`
	Href     string `json:"href"`
}

func (r YoutubeBot) DeferredAction(p *Payload) {
	text := strings.TrimSpace(p.Text)
	if text != "" {
		response := &IncomingWebhook{
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
			results := YouTubeVideoFeedResults{}
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

func (r YoutubeBot) Description() (description string) {
	// In addition to a Run method, each Robot must implement a Description method which
	// is just a simple string describing what the Robot does. This is used in the included
	// /c command which gives users a list of commands and descriptions
	return "Gets the most relevant YouTube result"
}
