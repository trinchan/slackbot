package robots

type RobotsBot struct {
}

func init() {
	RegisterRobot("c", new(RobotsBot))
}

func (r RobotsBot) Run(p *Payload) (slashCommandImmediateReturn string) {
	output := ""
	for command, r := range Robots {
		output = output + "\n" + command + " - " + r.Description() + "\n"
	}
	return output
}

func (r RobotsBot) Description() (description string) {
	return "Lists commands!\n\tUsage: You already know how to use this!\n\tExpected Response: This message!"
}
