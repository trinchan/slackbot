package robots

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type RollBot struct {
}

func init() {
	RegisterRobot("roll", new(RollBot))
}

func (roll RollBot) Run(p *Payload) (slashCommandImmediateReturn string) {
	go roll.DeferredAction(p)
	return ""
}

func (roll RollBot) DeferredAction(p *Payload) {
	response := &IncomingWebhook{
		Channel:     p.ChannelID,
		Username:    "Dice Bot",
		Text:        Roll(p),
		IconEmoji:   ":ghost:",
		UnfurlLinks: true,
		Parse:       ParseStyleFull,
	}

	MakeIncomingWebhookCall(response)
}

func (r RollBot) Description() (description string) {
	return "Roll an N-sided die!\n\tUsage: /roll {int}\n\tExpected Result: @user rolled an X out of Y!"
}

func Roll(p *Payload) (result string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	text := strings.TrimSpace(p.Text)
	var num int
	var err error
	if text == "" {
		num = 100
	} else {
		num, err = strconv.Atoi(text)
	}
	if err == nil && num > 0 {
		return fmt.Sprintf("@%s rolled a %d out of %d!", p.UserName, 1+r.Intn(num), num)
	} else {
		return fmt.Sprintf("@%s: Invalid input (%s): Must be integer greater than zero!", p.UserName, p.Text)
	}
}
