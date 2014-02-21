package robots

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
)

var Robots = make(map[string]func() Robot)
var Config = new(Configuration)
var ConfigDirectory = flag.String("c", ".", "Configuration directory (default .)")

func init() {
	flag.Parse()
	configFile := filepath.Join(*ConfigDirectory, "config.json")
	config, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal("Error opening config: ", err)
	}

	err = json.Unmarshal(config, Config)
	if err != nil {
		log.Fatal("Error parsing config: ", err)
	}
}

func RegisterRobot(command string, RobotInitFunction func() Robot) {
	if _, ok := Robots[command]; ok {
		log.Printf("There are two robots mapped to %s!", command)
	} else {
		log.Printf("Registered: %s", command)
		Robots[command] = RobotInitFunction
	}
}

func MakeIncomingWebhookCall(payload *IncomingWebhook) error {
	webhook := url.URL{
		Scheme: "https",
		Host:   Config.Domain + ".slack.com",
		Path:   "/services/hooks/incoming-webhook",
	}

	json_payload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	post_data := url.Values{}
	post_data.Set("payload", string(json_payload))
	post_data.Set("token", Config.Token)

	webhook.RawQuery = post_data.Encode()
	resp, err := http.PostForm(webhook.String(), post_data)
	if resp.StatusCode != 200 {
		message := fmt.Sprintf("ERROR: Non-200 Response from Slack Incoming Webhook API: %s", resp.Status)
		log.Println(message)
    }
	return err
}