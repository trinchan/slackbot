package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/schema"
	"github.com/trinchan/slackbot/robots"
)

func main() {
	http.HandleFunc("/slack", SlashCommandHandler)
	http.HandleFunc("/slack_hook", HookHandler)
	StartServer()
}

func HookHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	d := schema.NewDecoder()
	command := new(robots.OutgoingWebHook)
	err = d.Decode(command, r.PostForm)
	if err != nil {
		log.Println("Couldn't parse post request:", err)
	}
	log.Printf("Recieved command: %s from \"%s\"\n", command.Text[1:len(command.Text)], command.TeamDomain)
	//Assume outgoing webhooks are preceded with a character
	c := strings.Split(command.Text[1:len(command.Text)], " ")
	command.Robot = c[0]
	command.Text = strings.Join(c[1:], " ")

	robot := GetRobot(command.Robot)
	if robot == nil {
		jsonResp(w, "No robot for that command yet :(")
		return
	}
	w.WriteHeader(http.StatusOK)
	jsonResp(w, robot.Run(&command.Payload))
}

func SlashCommandHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	d := schema.NewDecoder()
	command := new(robots.SlashCommand)
	err = d.Decode(command, r.PostForm)
	if err != nil {
		log.Println("Couldn't parse post request:", err)
	}
	command.Robot = command.Command[1:]
	robot := GetRobot(command.Robot)
	if robot == nil {
		plainResp(w, "No robot for that command yet :(")
		return
	}
	plainResp(w, robot.Run(&command.Payload))
}

func jsonResp(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := map[string]string{"text": msg}
	r, err := json.Marshal(resp)
	if err != nil {
		log.Println("Couldn't marshal hook response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(r)
}

func plainResp(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(msg))
}

func StartServer() {
	port := robots.Config.Port
	log.Printf("Starting HTTP server on %d", port)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatal("Server start error: ", err)
	}
}

func GetRobot(command string) robots.Robot {
	if r, ok := robots.Robots[command]; ok {
		return r
	}
	return nil
}
