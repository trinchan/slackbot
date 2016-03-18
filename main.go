package main

import (
	_ "github.com/trinchan/slackbot/importer"
	"github.com/trinchan/slackbot/robots"
	"github.com/trinchan/slackbot/server"
)

func main() {
	server.Main(robots.Robots)
}
