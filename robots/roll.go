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
	RegisterRobot("roll", func() (robot Robot) { return new(RollBot) })
}

func (roll RollBot) Run(command *SlashCommand) (slashCommandImmediateReturn string) {
	go roll.DeferredAction(command)
	return ""
}

func (roll RollBot) DeferredAction(command *SlashCommand) {
	response := new(IncomingWebhook)
	response.Channel = command.Channel_ID
	response.Username = "Dice Bot"
	response.Text = Roll(command)
	response.Icon_Emoji = ":ghost:"
	response.Unfurl_Links = true
	response.Parse = "full"
	MakeIncomingWebhookCall(response)
}

func (r RollBot) Description() (description string) {
	return "Roll an N-sided die!\n\tUsage: /roll {int}\n\tExpected Result: @user rolled an X out of Y!"
}

func Roll(command *SlashCommand) (result string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	text := strings.TrimSpace(command.Text)
	var num int
	var err error
	if text == "" {
		num = 100
	} else {
		num, err = strconv.Atoi(command.Text)
	}
	if err == nil && num > 0 {
		return fmt.Sprintf("@%s rolled a %d out of %d!", command.User_Name, 1+r.Intn(num), num)
	} else {
		return fmt.Sprintf("@%s: Invalid input (%s): Must be integer greater than zero!", command.User_Name, command.Text)
	}
}
