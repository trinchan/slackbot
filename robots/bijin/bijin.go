package robots

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/trinchan/slackbot/robots"
)

type profile struct {
	Success bool          `json:"success"`
	Result  profileResult `json:"result"`
	Message string        `json:"message"`
}

type profileResult struct {
	ProfileInfo []profileInfo `json:"profile_info"`
}

type profileInfo struct {
	Title string      `json:"title"`
	Note  string      `json:"note"`
	URL   string      `json:"url"`
	Type  profileType `json:"type"`
}

type profileType int

const infoTempl = "http://www.bijint.com/assets/profile/%s/pc/ja/"
const picTempl = "http://www.bijint.com/assets/pict/%s/pc/"

var prefixes = map[string]string{
	"japan":     "jp",
	"thailand":  "thailand",
	"2012":      "2012jp",
	"2011":      "2011jp",
	"tokyo":     "tokyo",
	"hokkaido":  "hokkaido",
	"sendai":    "sendai",
	"akita":     "akita",
	"gunma":     "gunma",
	"niigata":   "niigata",
	"kanazawa":  "kanazawa",
	"fukui":     "fukui",
	"nagoya":    "nagoya",
	"kyoto":     "kyoto",
	"osaka":     "osaka",
	"kobe":      "kobe",
	"okayama":   "okayama",
	"kagawa":    "kagawa",
	"fukuoka":   "fukuoka",
	"kagoshima": "kagoshima",
	"okinawa":   "okinawa",
	"kumamoto":  "kumamoto",
	"saitama":   "saitama",
	"hiroshima": "hiroshima",
	"chiba":     "chiba",
	"nara":      "nara",
	"yamaguchi": "yamaguchi",
	"nagano":    "nagano",
	"shizuoka":  "shizuoka",
	"miyazaki":  "miyazaki",
	"tottori":   "tottori",
	"iwate":     "iwate",
	"ibaraki":   "ibaraki",
	"tochigi":   "tochigi",
	"taiwan":    "taiwan",
	"hawaii":    "hawaii",
	"seifuku":   "seifuku",
	"megane":    "megane",
	"sara":      "sara",
	"hairstyle": "hairstyle",
	"circuit":   "cc",
	"hanayome":  "hanayome",
	"waseda":    "wasedastyle",
}

type bot struct {
	Token    string
	Location *time.Location
}

func init() {
	b := &bot{}
	l, err := time.LoadLocation(os.Getenv("BIJIN_TIMEZONE"))
	if err != nil {
		log.Printf("Error loading location: %v", err)
	}
	b.Location = l
	b.Token = os.Getenv("BIJIN_SLACK_TOKEN")
	if b.Token == "" {
		log.Println("[WARNING] BIJIN_SLACK_TOKEN not set, will not verify payloads")
	}
	robots.RegisterRobot("bijin", b)
}

func (r bot) Run(p *robots.Payload) (slashCommandImmediateReturn string) {
	if r.Token != "" && p.Token != r.Token {
		log.Println("[WARNING] Ignoring request from unknown source")
		return ""
	}
	go r.DeferredAction(p)
	return ""
}

func (r bot) DeferredAction(p *robots.Payload) {
	response := &robots.IncomingWebhook{
		Domain:      p.TeamDomain,
		Channel:     p.ChannelID,
		Username:    "Bijin Bot",
		IconEmoji:   ":ghost:",
		UnfurlLinks: true,
	}

	rand.Seed(time.Now().UTC().UnixNano())
	t := time.Now().In(r.Location)
	hours := fmt.Sprintf("%02d", t.Hour())
	minutes := fmt.Sprintf("%02d", t.Minute())
	region, link, profileLink := getLink(strings.ToLower(strings.TrimSpace(p.Text)))
	prof := getProfile(profileLink, hours, minutes)
	response.Text = fmt.Sprintf("<@%s|%s> Here's your <%s%s%s.jpg|%s:%s 美人 (%s)>\n%s", p.UserID, p.UserName, link, hours, minutes, hours, minutes, strings.ToTitle(region), prof)
	response.Send()
}

func (p profile) String() string {
	if !p.Success {
		return ""
	}
	msg := ""
	if p.Message != "" {
		msg += p.Message
	}
	for _, r := range p.Result.ProfileInfo {
		if r.Title == "" || r.Note == "" || r.Note == "-" {
			continue
		}
		if r.URL != "" {
			msg += fmt.Sprintf("\n%s: <%s|%s>", r.Title, r.URL, r.Note)
		} else {
			msg += fmt.Sprintf("\n%s: %s", r.Title, r.Note)
		}
	}
	return msg
}

func getProfile(profileLink, hours, minutes string) profile {
	p := profile{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s%s.json", profileLink, hours, minutes), nil)
	if err != nil {
		return p
	}
	req.Header.Add("Host", "http://www.bijint.com")
	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return p
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return p
	}
	err = json.Unmarshal(b, &p)
	if err != nil {
		return p
	}
	return p
}

func getLink(region string) (string, string, string) {
	var r string
	var i string
	if i, ok := prefixes[region]; ok {
		return region, fmt.Sprintf(picTempl, i), fmt.Sprintf(infoTempl, i)
	}
	for r, i = range prefixes {
		break
	}
	return r, fmt.Sprintf(picTempl, i), fmt.Sprintf(infoTempl, i)
}

func (r bot) Description() (description string) {
	return "Displays the current time's 美人 (hope bijint.com doesn't get mad me)\n\tUsage: /bijin\n\tExpected Response: (beautiful woman)"
}
