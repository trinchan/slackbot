slackbot
===========
Simple, pluggable bot framework for [Slack](https://www.slack.com) chat.  

Dependencies
============
Schema  - `go get github.com/gorilla/schema`

Installation
============
You can grab the source code using `go get github.com/trinchan/slackbot` and install like usual. `go install github.com/trinchan/slackbot`

Setup
=====
###Setup an Incoming Webhook
If you don't already have an [Incoming Webhook](https://my.slack.com/services/new/incoming-webhook) setup for your bot, you'll want to start here.  Set one up (make sure it's enabled) and don't be afraid to read through the setup instructions.  You're after the "Webhook URL" that slack generates for you.  At the end of the URL, is the token that slack uses to authenticate where the URL is coming from:  
```
https://hooks.slack.com/services/AAAAAAAAA/BBBBBBBBB/YourSecretToken123456789
```
That's really all you care about right now.  You can set the default Icon, Name, and default Channel, but slack will let you override that information in the http requests you send.  So don't worry yourself with setting everything up.  You just need the token.

###Create a Domain Configuration File
Assuming you've already pulled the source, and successfully compiled/installed, you should have a `slackbot` executable in your `$GOPATH/bin`.  You need to create a file named `config.json` and give your bot the proper credentials to send messages to your slack server.  Feel free to place the file in the a sub-folder if you want to be all organized like that.  If you want to attach more than one slack server to your bot, you can simply add another entry under "domain_tokens".

The config file (config.json) has the following format:

```json
{
    "port": PORT_FOR_BOT,
    "domain_tokens": {
        "YOUR_SLACK_DOMAIN":       "YOUR_SLACK_INCOMING_WEBHOOK_TOKEN",
        "YOUR_OTHER_SLACK_DOMAIN": "MATCHING_INCOMING_WEBHOOK_TOKEN"
    }
}
```
Note that the last "domain_token" does NOT have a comma at the end of the line (but the others do)

###Send messages to your bot
This framework can respond to "slash commands" and "outgoing webhooks"  If you want users to be able to silently type `/ping`, and have the ping-bot respond in their channel, then you'll want to set up "slash commands".  Each bot will need it's own command setup.  The other option is to configure an outgoing webhook with a symbol for the prefix. Exe: `!ping`.  This option only requires one configuration, but the commands will be entered into the channel as regular messages.

#####Configuring an Outgoing Webhook
I use an [Outgoing Webhook](https://my.slack.com/services/new/outgoing-webhook)

1. Add a new Outgoing Webhook Integration.  
2. Here, you can tell slack to ONLY pay attention to a specific channel, or to simply listen to all public channels.  Outgoing Webhooks can not listen to private channels/direct messages.  
3. The {trigger_word} should only be one character (preferrably a symbol, such as ! or ?) and typing `{trigger_word}ping` will trigger the Ping bot.  
TODO: Clean up the trigger_word configuration.  Maybe something can be added to the config?
4. The URL should follow the following format: `your_address.com:port/slack_hook` (no trailing /)  
No other configuration is necessary.

The bot will respond to commands of the form `{trigger_word}bot param param param` in the specified channels
#####Configuring Slash Commands
Alternativly, each bot you make can respond to a corresponding [Slash Command](https://my.slack.com/services/new/slash-commands).

1. Add a new slash command, use the [bot's name](https://github.com/trinchan/slackbot/tree/master/robots) as the name of the command.  
2. The URL should follow the following format: `your_address.com:port/slack` (no trailing /)  
3. You want to use POST.  
4. This bot does not currently pay attention to the payload's token.  
TODO: Pay attention to the payload's token.
5. Repeat for each bot you want to enable.

The bot will respond to commands of the form `/bot param param param`

Adding Bots
===========
Create a new go file in the robots directory and, if necessary, a config file wherever you want to store the config files for your bots (all config files should be in the same directory and have names matching the bot name).

If you use [Sublime Text](http://www.sublimetext.com/) for development, then you can simply add and use the included snippet to easily generate a template Robot based on the filename. Otherwise, refer to the template below.

```go
package robots

import (
    "encoding/json"
    "flag"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"
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
    t := &TestBot{}
    RegisterRobot("test", t)
}

// All Robots must implement a Run command to be executed when the registered command is received.
func (r TestBot) Run(p *Payload) string {
    // If you (optionally) want to do some asynchronous work (like sending API calls to slack)
    // you can put it in a go routine like this
    go r.DeferredAction(p)
    // The string returned here will be shown only to the user who executed the command
    // and will show up as a message from slackbot.
    return "Text to be returned only to the user who made the command."
}

func (r TestBot) DeferredAction(p *Payload) {
    // Let's use the IncomingWebhook struct defined in definitions.go to form and send an
    // IncomingWebhook message to slack that can be seen by everyone in the room. You can
    // read the Slack API Docs (https://api.slack.com/) to know which fields are required, etc.
    // You can also see what data is available from the command structure in definitions.go
    response := &IncomingWebhook{
        Domain:      p.TeamDomain
        Channel:     p.ChannelID,
        Username:    "Test Bot",
        Text:        "Hi there!",
        IconEmoji:   ":ghost:",
        UnfurlLinks: true,
        Parse:       ParseStyleFull,
    }
    response.Send()
}

func (r TestBot) Description() (description string) {
    // In addition to a Run method, each Robot must implement a Description method which
    // is just a simple string describing what the Robot does. This is used in the included
    // /c command which gives users a list of commands and descriptions
    return "This is a description for TestBot which will be displayed on /c"
}

```

If you are using [Slash Commands](https://my.slack.com/services/new/slash-commands) instead of an outgoing-webhook, you'll need to add a new slash command integration for each bot you add.

Running
=======
`slackbot -c {PATH_TO_CONFIG_FILE_DIRECTORY}`

If you see output similar to below and you have the commands enabled in your Slack integration, you're ready to go!
```
2015/05/03 13:39:48 WARNING: Could not find configuration file bijin.json in slack
2015/05/03 13:39:48 Registered: bijin
2015/05/03 13:39:48 Registered: decide
2015/05/03 13:39:48 Registered: ping
2015/05/03 13:39:48 WARNING: Could not find configuration file pivotal.json in slack
2015/05/03 13:39:48 Registered: pivotal
2015/05/03 13:39:48 Registered: c
2015/05/03 13:39:48 Registered: roll
2015/05/03 13:39:48 Found 4 domain configurations
2015/05/03 13:39:48 Port: 8888
2015/05/03 13:39:48 Registered: store
2015/05/03 13:39:48 Registered: wiki
2015/05/03 13:39:48 Registered: youtube
2015/05/03 13:39:48 Starting HTTP server on 8888
```
