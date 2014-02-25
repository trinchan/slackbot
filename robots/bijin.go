package robots

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type BijinBot struct {
}

type BijinConfiguration struct {
	Timezone string `json:"timezone"`
}

var BijinConfig = new(BijinConfiguration)

func init() {
	flag.Parse()
	configFile := filepath.Join(*ConfigDirectory, "bijin.json")
	if _, err := os.Stat(configFile); err == nil {
		config, err := ioutil.ReadFile(configFile)
		if err != nil {
			log.Printf("ERROR: Error opening bijin config: %s", err)
			return
		}
		err = json.Unmarshal(config, BijinConfig)
		if err != nil {
			log.Printf("ERROR: Error parsing bijin config: %s", err)
			return
		}
	} else {
		log.Printf("WARNING: Could not find configuration file bijin.json in %s", *ConfigDirectory)
	}
	RegisterRobot("/bijin", func() (robot Robot) { return new(BijinBot) })
}

func (r BijinBot) Run(command *SlashCommand) (slashCommandImmediateReturn string) {
	go r.DeferredAction(command)
	return ""
}

func (r BijinBot) DeferredAction(command *SlashCommand) {
	response := new(IncomingWebhook)
	response.Channel = command.Channel_ID
	response.Username = "Bijin Bot"
	response.Icon_Emoji = ":ghost:"
	response.Unfurl_Links = true

	rand.Seed(time.Now().UTC().UnixNano())
	var t time.Time
	if BijinConfig.Timezone != "" {
		loc, err := time.LoadLocation(BijinConfig.Timezone)
		if err == nil {
			t = time.Now().In(loc)
		} else {
			t = time.Now()
			response.Text = fmt.Sprintf("ERROR: Unknown timezone (%s) - Serving UTC", BijinConfig.Timezone)
			MakeIncomingWebhookCall(response)
		}
	} else {
		t = time.Now()
		response.Text = fmt.Sprintf("WARNING: No timezone set - Serving UTC")
		MakeIncomingWebhookCall(response)
	}
	hours := fmt.Sprintf("%02d", t.Hour())
	minutes := fmt.Sprintf("%02d", t.Minute())
	links := []string{"http://www.bijint.com/jp/tokei_images/", "http://www.bijint.com/jp/img/clk/"}
	response.Text = fmt.Sprintf("<@%s|%s> Here's your <%s%s%s.jpg|%s:%s 美人>", command.User_ID, command.User_Name, links[rand.Intn(len(links))], hours, minutes, hours, minutes)
	MakeIncomingWebhookCall(response)
}

func (r BijinBot) Description() (description string) {
	return "Displays the current time's 美人 (hope bijint.com doesn't get mad me)\n\tUsage: /bijin\n\tExpected Response: (beautiful woman)"
}
