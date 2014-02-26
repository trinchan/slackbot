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
	"strings"
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
	t := time.Now()
	if BijinConfig.Timezone != "" {
		loc, err := time.LoadLocation(BijinConfig.Timezone)
		if err == nil {
			t = t.In(loc)
		} else {
			response.Text = fmt.Sprintf("ERROR: Unknown timezone (%s) - Serving UTC", BijinConfig.Timezone)
			MakeIncomingWebhookCall(response)
		}
	} else {
		response.Text = fmt.Sprintf("WARNING: No timezone set - Serving UTC")
		MakeIncomingWebhookCall(response)
	}
	hours := fmt.Sprintf("%02d", t.Hour())
	minutes := fmt.Sprintf("%02d", t.Minute())
	region, link := GetLink(strings.ToLower(strings.TrimSpace(command.Text)))
	response.Text = fmt.Sprintf("<@%s|%s> Here's your (%s) <%s%s%s.jpg|%s:%s 美人> ", command.User_ID, command.User_Name, strings.ToTitle(region), link, hours, minutes, hours, minutes)
	MakeIncomingWebhookCall(response)
}

func GetLink(region string) (string, string) {
	links := map[string]string{"japan": "http://www.bijint.com/jp/tokei_images/",
		"clk":       "http://www.bijint.com/jp/img/clk/",
		"2012":      "http://www.bijint.com/2012jp/tokei_images/",
		"2011":      "http://www.bijint.com/2011jp/tokei_images/",
		"tokyo":     "http://www.bijint.com/tokyo/tokei_images/",
		"hokkaido":  "http://www.bijint.com/hokkaido/tokei_images/",
		"sendai":    "http://www.bijint.com/sendai/tokei_images/",
		"akita":     "http://www.bijint.com/akita/tokei_images/",
		"gunma":     "http://www.bijint.com/gunma/tokei_images/",
		"niigata":   "http://www.bijint.com/niiagata/tokei_images/",
		"kanazawa":  "http://www.bijint.com/kanazawa/tokei_images/",
		"fukui":     "http://www.bijint.com/fukui/tokei_images/",
		"nagoya":    "http://www.bijint.com/nagoya/tokei_images/",
		"kyoto":     "http://www.bijint.com/kyoto/tokei_images/",
		"osaka":     "http://www.bijint.com/osaka/tokei_images/",
		"kobe":      "http://www.bijint.com/kobe/tokei_images/",
		"okayama":   "http://www.bijint.com/okayama/tokei_images/",
		"kagawa":    "http://www.bijint.com/kagawa/tokei_images/",
		"fukuoka":   "http://www.bijint.com/fukuoka/tokei_images/",
		"kagoshima": "http://www.bijint.com/kagoshima/tokei_images/",
		"okinawa":   "http://www.bijint.com/okinawa/tokei_images/",
		"kumamoto":  "http://www.bijint.com/kumamoto/tokei_images/",
		"saitama":   "http://www.bijint.com/saitama/tokei_images/",
		"hiroshima": "http://www.bijint.com/hiroshima/tokei_images/",
		"chiba":     "http://www.bijint.com/chiba/tokei_images/",
		"nara":      "http://www.bijint.com/nara/tokei_images/",
		"yamaguchi": "http://www.bijint.com/yamaguchi/tokei_images/",
		"nagano":    "http://www.bijint.com/nagano/tokei_images/",
		"shizuoka":  "http://www.bijint.com/shizuoka/tokei_images/",
		"miyazaki":  "http://www.bijint.com/miyazaki/tokei_images/",
		"tottori":   "http://www.bijint.com/tottori/tokei_images/",
		"iwate":     "http://www.bijint.com/iwate/tokei_images/",
		"ibaraki":   "http://www.bijint.com/ibaraki/tokei_images/",
		"tochigi":   "http://www.bijint.com/tochigi/tokei_images/",
		"taiwan":    "http://www.bijint.com/taiwan/tokei_images/",
		"hawaii":    "http://www.bijint.com/hawaii/tokei_images/",
		"seifuku":   "http://www.bijint.com/seifuku/tokei_images/",
		"megane":    "http://www.bijint.com/megane/tokei_images/",
		"sara":      "http://www.bijint.com/sara/tokei_images/",
		"hairstyle": "http://www.bijint.com/hairstyle/tokei_images/",
		"asahi":     "http://www.bijint.com/tv-asahi/tokei_images/",
		"circuit":   "http://www.bijint.com/cc/tokei_images/",
		"hanayome":  "http://www.bijint.com/hanayome/tokei_images/",
		"waseda":    "http://www.bijint.com/waseda/tokei_images/",
	}
	if link, ok := links[region]; ok {
		return region, link
	} else {
		for region, link := range links {
			return region, link
		}
	}
	return "japan", links["japan"]
}

func (r BijinBot) Description() (description string) {
	return "Displays the current time's 美人 (hope bijint.com doesn't get mad me)\n\tUsage: /bijin\n\tExpected Response: (beautiful woman)"
}
