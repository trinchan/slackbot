package robots
import (
    "encoding/json"
    "flag"
    "path/filepath"
    "io/ioutil"
    "log"
    "os"
)
type StoreBot struct {
}

type StoreConfiguration struct {
}

var StoreConfig = new(StoreConfiguration)

// Loads the config file and registers the bot with the server for command /store.
func init() {
    flag.Parse()
    configFile := filepath.Join(*ConfigDirectory, "store.json")
    if _, err := os.Stat(configFile); err == nil {
        config, err := ioutil.ReadFile(configFile)
        if err != nil {
            log.Printf("ERROR: Error opening store config: %s", err)
            return
        }
        err = json.Unmarshal(config, StoreConfig)
        if err != nil {
            log.Printf("ERROR: Error parsing store config: %s", err)
            return
        }
    } else {
        log.Printf("WARNING: Could not find configuration file store.json in %s", *ConfigDirectory)
    }
    RegisterRobot("/store", func() (robot Robot) { return new(StoreBot) })
}

// All Robots must implement a Run command to be executed when the registered command is received.
func (r StoreBot) Run(command *SlashCommand) (slashCommandImmediateReturn string) {
    // If you (optionally) want to do some asynchronous work (like sending API calls to slack)
    // you can put it in a go routine like this
    go r.DeferredAction(command)
    // The string returned here will be shown only to the user who executed the command
    // and will show up as a message from slackbot.
    return ""
}

func (r StoreBot) DeferredAction(command *SlashCommand) {
    // Let's use the IncomingWebhook struct defined in definitions.go to form and send an
    // IncomingWebhook message to slack that can be seen by everyone in the room. You can
    // read the Slack API Docs (https://api.slack.com/) to know which fields are required, etc.
    // You can also see what data is available from the command structure in definitions.go
    response := new(IncomingWebhook)
    response.Channel = command.Channel_ID
    response.Username = "Store Bot"
    response.Text = fmt.Sprintf(":famima: @group @%s wants to go to the store :break:", command.User_Name)
    response.Icon_Emoji = ":famima:"
    response.Unfurl_Links = true
    response.Parse = "full"
    MakeIncomingWebhookCall(response)
}

func (r StoreBot) Description() (description string) {
    // In addition to a Run method, each Robot must implement a Description method which
    // is just a simple string describing what the Robot does. This is used in the included
    // /c command which gives users a list of commands and descriptions
    return "This is a description for StoreBot which will be displayed on /c"
}