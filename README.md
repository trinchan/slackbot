slackbot
===========

A framework for building slash commands for Slack in Go

Dependencies
============
Schema  - `go get github.com/gorilla/schema`

Installation
============
You can grab the source code using `go get github.com/trinchan/slackbot` and install like usual.

Setup
=====
Create a config file (config.json) with the following format:

```json
{
    "domain": "{YOUR_SLACK_DOMAIN}",
    "port": {PORT_FOR_BOT},
    "token": "{YOUR_SLACK_INCOMING_WEBHOOK_TOKEN}"
}
```

Make sure you have [Incoming Webhooks](https://my.slack.com/services/new/incoming-webhook) enabled and you are using that integration token for your config.

For each slash command (including the default commands!), be sure to add a corresponding entry in [Slash Commands](https://my.slack.com/services/new/slash-commands) to POST to server:port/slack of your slackbot setup. Note no trailing slash after /slack.

Adding Bots
===========
Create a new go file in the robots directory and, if necessary, a config file wherever you want to store the config files for your bots (all config files should be in the same directory and have names matching the bot name).

If you use [Sublime Text](http://www.sublimetext.com/) for development, then you can simply add and use the included snippet to easily generate a template Robot based on the filename. Otherwise, refer to the template below.

```go
package robots
import (
    "encoding/json"
    "flag"
    "path/filepath"
    "io/ioutil"
    "log"
    "os"
)
type TestBot struct {
}

type TestConfiguration struct {
}

var TestConfig = new(TestConfiguration)

// Loads the config file and registers the bot with the server for command /test.
func init() {
    flag.Parse()
    configFile := filepath.Join(*ConfigDirectory, "test.json")
    if _, err := os.Stat(configFile); err == nil {
        config, err := ioutil.ReadFile(configFile)
        if err != nil {
            log.Printf("ERROR: Error opening test config: %s", err)
            return
        }
        err = json.Unmarshal(config, TestConfig)
        if err != nil {
            log.Printf("ERROR: Error parsing test config: %s", err)
            return
        }
    } else {
        log.Printf("WARNING: Could not find configuration file test.json in %s", *ConfigDirectory)
    }
    RegisterRobot("/test", func() (robot Robot) { return new(TestBot) })
}

// All Robots must implement a Run command to be executed when the registered command is received.
func (r TestBot) Run(command *SlashCommand) (slashCommandImmediateReturn string) {
    // If you (optionally) want to do some asynchronous work (like sending API calls to slack)
    // you can put it in a go routine like this
    go r.DeferredAction(command)
    // The string returned here will be shown only to the user who executed the command
    // and will show up as a message from slackbot.
    return "Text to be returned only to the user who made the command."
}

func (r TestBot) DeferredAction(command *SlashCommand) {
    // Let's use the IncomingWebhook struct defined in definitions.go to form and send an
    // IncomingWebhook message to slack that can be seen by everyone in the room. You can
    // read the Slack API Docs (https://api.slack.com/) to know which fields are required, etc.
    // You can also see what data is available from the command structure in definitions.go
    response := new(IncomingWebhook)
    response.Channel = command.Channel_ID
    response.Username = "Test Bot"
    response.Text = "Hi there!"
    response.Icon_Emoji = ":ghost:"
    response.Unfurl_Links = true
    response.Parse = "full"
    MakeIncomingWebhookCall(response)
}

func (r TestBot) Description() (description string) {
    // In addition to a Run method, each Robot must implement a Description method which
    // is just a simple string describing what the Robot does. This is used in the included
    // /c command which gives users a list of commands and descriptions
    return "This is a description for TestBot which will be displayed on /c"
}
```

Now just add a corresponding entry in [Slash Commands](https://my.slack.com/services/new/slash-commands) to POST to server:port/slack of your slackbot setup. Note no trailing slash after /slack.

Running
=======
`slackbot -c {PATH_TO_CONFIG_FILE_DIRECTORY}`

If you see output similar to below and you have the commands enabled in your Slack integration, you're ready to go!
```
2014/02/18 10:55:07 Registered: /decide
2014/02/18 10:55:07 Registered: /ping
2014/02/18 10:55:07 Registered: /c
2014/02/18 10:55:07 Registered: /roll
2014/02/18 10:55:07 Starting HTTP server on 8888
```
