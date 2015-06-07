package robots

import "log"

// Robot describes the necessary methods to be registered as a slack bot
type Robot interface {
	Run(p *Payload) (botString string)
	Description() (description string)
}

// Robots is the map of registered command to robot
var Robots = make(map[string][]Robot)

// RegisterRobot registers a robot in the Robots map with
func RegisterRobot(command string, r Robot) {
	log.Printf("Registered: %s", command)
	Robots[command] = append(Robots[command], r)
}
