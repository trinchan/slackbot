#!/bin/bash
touch $1
> $1
robots=(
    "github.com/trinchan/slackbot/robots/decide"
    "github.com/trinchan/slackbot/robots/bijin"
    "github.com/trinchan/slackbot/robots/nihongo"
    "github.com/trinchan/slackbot/robots/ping"
    "github.com/trinchan/slackbot/robots/roll"
    "github.com/trinchan/slackbot/robots/store"
    "github.com/trinchan/slackbot/robots/wiki"
    "github.com/trinchan/slackbot/robots/bot"
)

echo "package importer

import (" >> $1

for robot in "${robots[@]}"
do
    echo "    _ \"$robot\" // automatically generated import to register bot, do not change" >> $1
done
echo ")" >> $1

gofmt -w -s $1
