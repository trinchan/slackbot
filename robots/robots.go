package robots

type RobotsBot struct {
}

func init() {
	RegisterRobot("/c", func() (robot Robot) { return new(RobotsBot) })
}

func (r RobotsBot) Run(command *SlashCommand) (slashCommandImmediateReturn string) {
	output := ""
	for command, RobotInitFunction := range Robots {
		robot := RobotInitFunction()
		output = output + "\n" + command + " - " + robot.Description() + "\n"
	}
	return output
}

func (r RobotsBot) DeferredAction(command *SlashCommand) {
}

func (r RobotsBot) Description() (description string) {
	return "Lists commands!\n\tUsage: You already know how to use this!\n\tExpected Response: This message!"
}
