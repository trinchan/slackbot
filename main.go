package main

import (
    "github.com/gorilla/schema"
    "github.com/trinchan/slackbot/robots"
    "io"
    "encoding/json"
    "log"
    "net/http"
    "strconv"
    "strings"
)

func main() {
    http.HandleFunc("/slack", CommandHandler)
    http.HandleFunc("/slack_hook", CommandHandler)
    StartServer()
}
func CommandHandler(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    hook := r.URL.Path == "/slack_hook"
    if err == nil {
        decoder := schema.NewDecoder()
        command := new(robots.SlashCommand)
        err := decoder.Decode(command, r.PostForm)
        if err != nil {
            log.Println("Couldn't parse post request:", err)
        }
        if hook {
            c := strings.Split(command.Text, " ")
            command.Command = "/" + c[1]
            command.Text = strings.Join(c[2:], " ")
        }
        robot := GetRobot(command)
        w.WriteHeader(http.StatusOK)
        if robot != nil {
            if hook {
                JSONResp(w, robot.Run(command))
            } else {
                plainResp(w, robot.Run(command))
            }
        } else {
            r := "No robot for that command yet :("
            if hook {
                JSONResp(w, r)
            } else {
                plainResp(w, r)
            }
        }
    }
}

func JSONResp(w http.ResponseWriter, msg string) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    resp := map[string]string{"text": msg}
    r, err := json.Marshal(resp)
    if err != nil {
        log.Println("Couldn't marshal hook response:", err)
    } else {
        io.WriteString(w, string(r))
    }
}

func plainResp(w http.ResponseWriter, msg string) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    io.WriteString(w, msg)
}

func StartServer() {
    port := robots.Config.Port
    log.Printf("Starting HTTP server on %d", port)
    err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
    if err != nil {
        log.Fatal("Server start error: ", err)
    }
}

func GetRobot(command *robots.SlashCommand) (robot robots.Robot) {
    if RobotInitFunction, ok := robots.Robots[command.Command]; ok {
        return RobotInitFunction()
    } else {
        return nil
    }
}
