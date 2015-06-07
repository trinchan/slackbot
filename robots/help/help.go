package directory

import "github.com/trinchan/slackbot/robots"

type bot struct{}

func init() {
	r := &bot{}
	robots.RegisterRobot("help", r)
}

func (r *bot) Run(p *robots.Payload) (slashCommandImmediateReturn string) {
	output := ""
	for command, robs := range robots.Robots {
		for _, r := range robs {
			output = output + "\n" + command + " - " + r.Description() + "\n"
		}
	}
	return output
}

func (r *bot) Description() (description string) {
	return "Lists commands!\n\tUsage: You already know how to use this!\n\tExpected Response: This message!"
}
