package robots

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/trinchan/slackbot/robots"
)

type bot struct {
	Random *rand.Rand
}

func init() {
	r := &bot{}
	r.Random = rand.New(rand.NewSource(time.Now().UnixNano()))

	robots.RegisterRobot("roll", r)
}

func (r bot) Run(p *robots.Payload) (slashCommandImmediateReturn string) {
	go r.DeferredAction(p)
	return ""
}

func (r bot) DeferredAction(p *robots.Payload) {
	response := &robots.IncomingWebhook{
		Domain:      p.TeamDomain,
		Channel:     p.ChannelID,
		Username:    "Dice Bot",
		Text:        r.roll(p),
		IconEmoji:   ":ghost:",
		UnfurlLinks: true,
		Parse:       robots.ParseStyleFull,
	}

	response.Send()
}

func (r bot) Description() (description string) {
	return "Roll an N-sided die!\n\tUsage: /roll {int}\n\tExpected Result: @user rolled an X out of Y!"
}

func (r bot) roll(p *robots.Payload) (result string) {
	text := strings.TrimSpace(p.Text)
	var num int
	var err error
	if text == "" {
		num = 100
	} else {
		num, err = strconv.Atoi(text)
	}
	if err == nil && num > 0 {
		return fmt.Sprintf("@%s rolled a %d out of %d!", p.UserName, 1+r.Random.Intn(num), num)
	}
	return fmt.Sprintf("@%s: Invalid input (%s): Must be integer greater than zero!", p.UserName, p.Text)
}
