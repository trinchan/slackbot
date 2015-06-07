package bot

import (
	"fmt"
	"strings"

	"github.com/trinchan/slackbot/robots"
)

type bot struct{}

func init() {
	s := &bot{}
	robots.RegisterRobot("bot", s)
}

func (r bot) Run(p *robots.Payload) (slashCommandImmediateReturn string) {
	text := strings.TrimSpace(p.Text)
	resp := ""
	if text != "" {
		com := strings.Split(text, " ")
		r := com[0]
		robs := robots.Robots[r]
		if len(robs) == 0 {
			return fmt.Sprintf("%s bot not found", r)
		}
		fp := &robots.Payload{}
		*fp = *p
		fp.Text = strings.Join(com[1:], " ")
		fp.Robot = r
		for _, rob := range robs {
			resp += fmt.Sprintf("\n%s", rob.Run(fp))
		}
	}
	return resp
}

func (r bot) Description() (description string) {
	return "The bot command is a helper bot that can invoke other bots. Useful if you are integration limited."
}
