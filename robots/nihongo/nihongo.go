package nihongo

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
	w := &bot{}
	robots.RegisterRobot("nihongo", w)
}

type nihongoResponse struct {
	Search  string  `json:"search"`
	Entries []entry `json:"entries"`
}

type entry struct {
	Word       string `json:"word"`
	Furigana   string `json:"furigana"`
	Definition string `json:"definition"`
	Common     common `json:"common"`
}

type common bool

func (c common) String() string {
	if c {
		return "[COMMON]"
	}
	return ""
}

func (n bot) Run(p *robots.Payload) (slashCommandImmediateReturn string) {
	text := strings.TrimSpace(p.Text)
	if text != "" {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://nihongo.io/search?text=%s", url.QueryEscape(text)), nil)
		if err != nil {
			log.Print(err)
			return "Error getting definition from Nihongo.io :("
		}
		req.Header.Add("Accept", "application/json")
		resp, err := http.DefaultClient.Do(req)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			message := fmt.Sprintf("ERROR: Non-200 Response from Nihongo.io: %s", resp.Status)
			log.Println(message)
			return message
		} else if err != nil {
			return "Error getting definition from Nihongo.io :("
		} else {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Print(err)
				return "Error getting definition from Nihongo.io :("
			}
			r := nihongoResponse{}
			err = json.Unmarshal(b, &r)
			if err != nil {
				log.Print(err)
				return "Error getting definition from Nihongo.io :("
			}
			msg := fmt.Sprintf("\nSearch term: %s", r.Search)
			for i, entry := range r.Entries {
				msg += fmt.Sprintf("\n%d. %s %s (%s) - %s", i+1, entry.Common, entry.Word, entry.Furigana, entry.Definition)
			}
			return msg
		}
	}
	return ""
}
func (n bot) Description() (description string) {
	return "Nihongo bot!\n\tUsage: /nihongo {word}\n\tExpected Response: @user: {definition} "
}
