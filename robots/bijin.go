package robots
import (
    "encoding/json"
    "flag"
    "fmt"
    "math/rand"
    "path/filepath"
    "io/ioutil"
    "log"
    "os"
    "time"
)
type BijinBot struct {
}

type BijinConfiguration struct {
    Timezone string `json:"timezone"`
}

var BijinConfig = new(BijinConfiguration)

// Loads the config file and registers the bot with the server for command /bijint.
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

// All Robots must implement a Run command to be executed when the registered command is received.
func (r BijinBot) Run(command *SlashCommand) (slashCommandImmediateReturn string) {
    // If you (optionally) want to do some asynchronous work (like sending API calls to slack)
    // you can put it in a go routine like this
    go r.DeferredAction(command)
    // The string returned here will be shown only to the user who executed the command
    // and will show up as a message from slackbot.
    return ""
}

func (r BijinBot) DeferredAction(command *SlashCommand) {
    // Let's use the IncomingWebhook struct defined in definitions.go to form and send an
    // IncomingWebhook message to slack that can be seen by everyone in the room. You can
    // read the Slack API Docs (https://api.slack.com/) to know which fields are required, etc.
    // You can also see what data is available from the command structure in definitions.go
    response := new(IncomingWebhook)
    response.Channel = command.Channel_ID
    response.Username = "Bijin Bot"
    rand.Seed(time.Now().UTC().UnixNano())
    loc, err := time.LoadLocation(BijinConfig.Timezone)
    t := time.Now().In(loc)
    hours := fmt.Sprintf("%02d", t.Hour())
    minutes := fmt.Sprintf("%02d", t.Minute())
    links := []string{"http://www.bijint.com/jp/tokei_images/", "http://www.bijint.com/jp/img/clk/"}
    if err == nil {
        response.Text = fmt.Sprintf("<@%s|%s> Here's your <%s%s%s.jpg|%s:%s 美人>", command.User_ID, command.User_Name, links[rand.Intn(len(links))], hours, minutes, hours, minutes)
    } else {
        response.Text = fmt.Sprintf("ERROR")
    }

    response.Icon_Emoji = ":ghost:"
    response.Unfurl_Links = true
    // response.Parse = "full"
    MakeIncomingWebhookCall(response)
}
func (r BijinBot) Description() (description string) {
    // In addition to a Run method, each Robot must implement a Description method which
    // is just a simple string describing what the Robot does. This is used in the included
    // /c command which gives users a list of commands and descriptions
    return "Displays the current time's 美人 (hope bijint.com doesn't get mad me)\n\tUsage: /bijin\n\tExpected Response: (beautiful woman)"
}