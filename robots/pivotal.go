package robots
import (
    "encoding/json"
    "fmt"
    "flag"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "os"
    "path/filepath"
    "strings"
)

type PivotalBot struct {
}

type PivotalConfiguration struct {
    Project_ID int `json:"project_id"`
    Token   string    `json:"token"`
}
type SearchResultContainer struct {
    Epics EpicsSearchResult `json:epics`
    Query string `json:"query"`
    Stories StoriesSearchResult `json:stories`
}

type EpicsSearchResult struct {
    Total_Hits int `json:"total_hits"`
    Epics []Epic `json:"epics"`
    Total_Hits_With_Done int `json:"total_hits_with_done,omitempty"`
}

type StoriesSearchResult struct {
    Total_Hits int `json:"total_hits"`
    Total_Hits_With_Done int `json:"total_hits_with_done,omitempty"`
    Stories []Story `json:"stories"`
    Total_Points int `json:"total_points"`
    Total_Points_Completed int `json:"total_points_completed,omitempty"`
}

type Label struct {
    Kind string `json:"kind"`
    Created_At string `json:"created_at"`
    Updated_At string `json:"updated_at"`
    Name string `json:"name"`
    Project_ID int `json:"project_id"`
    ID int `json:"id"`
}

type Epic struct {
    Comment_IDs []int `json:"comment_ids"`
    Name string `json:"name"`
    Kind string `json:"kind"`
    ID int `json:"id"`
    Before_ID int `json:"before_id"`
    Updated_At string `json:"updated_at"`
    URL string `json:"url"`
    Project_ID int `json:"project_id"`
    Label_ID int `json:"label_id"`
    Created_At string `json:"created_at"`
    Description string `json:"description"`
    After_ID int `json:"after_id"`
}

type Story struct {
    Comment_IDs []int `json:"comment_ids,omitempty"`
    Current_State string `json:"current_state"`
    Deadline string `json:"deadline,omitempty"`
    Requested_By_ID int `json:"requested_by_id,omitempty"`
    Integration_ID int `json:"integration_id,omitempty"`
    Name string `json:"name"`
    Owned_By_ID int `json:"owned_by_id,omitempty"`
    Kind string `json:"kind"`
    Labels []Label `json:"labels,omitempty"`
    ID int `json:"id"`
    Planned_Iteration_Number int `json:"planned_iteration_number,omitempty"`
    External_ID string `json:"external_id,omitempty"`
    Estimate int `json:"estimate"`
    Task_IDs []int `json:"task_ids,omitempty"`
    Updated_At string `json:"updated_at"`
    URL string `json:"url"`
    Project_ID int `json:"project_id"`
    Story_Type string `json:"story_type"`
    Accepted_At string `json:"accepted_at"`
    Follower_IDs []int `json:"follower_ids,omitempty"`
    Created_At string `json:"created_at"`
    Description string `json:"description"`
    Owner_IDs []int `json:"owner_ids,omitempty"`
}

type Task struct {
    Complete bool `json:"complete"`
    Kind string `json:"kind"`
    ID int `json:"id"`
    Position int `json:"position"`
    Updated_At string `json:"updated_at"`
    Story_ID int `json:"story_id"`
    Created_At string `json:"created_at"`
    Description string `json:"description"`
}

var PivotalConfig = new(PivotalConfiguration)

// Loads the config file and registers the bot with the server for command /${1/(.+)/\L\1/g}.
func init() {
    flag.Parse()
    configFile := filepath.Join(*ConfigDirectory, "pivotal.json")
    if _, err := os.Stat(configFile); err == nil {
        config, err := ioutil.ReadFile(configFile)
        if err != nil {
            log.Printf("ERROR: Error opening pivotal config: %s", err)
            return
        }
        err = json.Unmarshal(config, PivotalConfig)
        if err != nil {
            log.Printf("ERROR: Error parsing pivotal config: %s", err)
            return
        }
    } else {
        log.Printf("WARNING: Could not find configuration file pivotal.json in %s", *ConfigDirectory)
    }
    RegisterRobot("/pivotal", func() (robot Robot) { return new(PivotalBot) })
}

// All Robots must implement a Run command to be executed when the registered command is received.
func (r PivotalBot) Run(command *SlashCommand) (slashCommandImmediateReturn string) {
    // The string returned here will be shown only to the user who executed the command
    // and will show up as a message from slackbot.
    text := strings.TrimSpace(command.Text)
    if text != "" {
        split := strings.Split(text, " ")
        pivotal_command := split[0]
        query := strings.Join(split[1:], " ")
        switch pivotal_command {
            case "query":
                return r.Query(query)
            case "start", "unstart", "unschedule", "finish", "accept", "reject", "deliver":
                return r.ChangeState(pivotal_command, query)
        }
        return fmt.Sprintf("Unknown pivotal command: %s\n%s", pivotal_command, r.Description())
    } else {
        return ""
    }
}

func (r PivotalBot) ChangeState(new_state string, story_id string) (result string) {
        put_parameters := url.Values{}
        put_parameters.Set("current_state", new_state + "ed")
        req, err := http.NewRequest("PUT", fmt.Sprintf("https://www.pivotaltracker.com/services/v5/projects/%d/stories/%s", PivotalConfig.Project_ID, story_id), nil)
        if err != nil {
            return fmt.Sprintf("Error forming put request to Pivotal: %s", err)
        }
        req.URL.RawQuery = put_parameters.Encode()
        req.Header.Add("X-TrackerToken", PivotalConfig.Token)
        resp, err := http.DefaultClient.Do(req)
        if err != nil {
            return fmt.Sprintf("Error making put request to Pivotal: %s", err)
        }
        defer resp.Body.Close()
        contents, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            return fmt.Sprintf("Error reading response body from Pivotal: %s", err)
        }
        story := new(Story)
        err = json.Unmarshal(contents, &story)
        if err != nil {
            return fmt.Sprintf("Couldn't unmarshal pivotal story response into struct: %s", err)
        }
        return fmt.Sprintf("[%s <%s|#%d>] - %s", story.Current_State, story.URL, story.ID, story.Name)
}

func (r PivotalBot) Query(query string) (result string) {
        get_parameters := url.Values{}
        get_parameters.Set("query", query)
        req, err := http.NewRequest("GET", fmt.Sprintf("https://www.pivotaltracker.com/services/v5/projects/%d/search", PivotalConfig.Project_ID), nil)
        req.URL.RawQuery = get_parameters.Encode()
        if err != nil {
            return fmt.Sprintf("Error forming get request to Pivotal: %s", err)
        }
        req.Header.Add("X-TrackerToken", PivotalConfig.Token)
        resp, err := http.DefaultClient.Do(req)
        if err != nil {
            return fmt.Sprintf("Error making get request to Pivotal: %s", err)
        }
        defer resp.Body.Close()
        contents, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            return fmt.Sprintf("Error reading response body from Pivotal: %s", err)
        }
        searchResults := new(SearchResultContainer)
        err = json.Unmarshal(contents, &searchResults)
        if err != nil {
            return fmt.Sprintf("Couldn't unmarshal pivotal response into struct: %s", err)
        }
        output := ""
        if searchResults.Stories.Total_Hits > 0 {
            for _, story := range searchResults.Stories.Stories {
                output = output + (fmt.Sprintf("[%s <https://www.pivotaltracker.com/s/projects/%d/stories/%d|#%d>] - %s", story.Current_State, story.Project_ID, story.ID, story.ID, story.Name)) + "\n"
            }
        }
        if output == "" {
            return "No stories matching that query :("
        } else {
            return strings.TrimSpace(output)
        }
}

func (r PivotalBot) Description() (description string) {
    // In addition to a Run method, each Robot must implement a Description method which
    // is just a simple string describing what the Robot does. This is used in the included
    // /c command which gives users a list of commands and descriptions
    return "Interact with the Pivotal API!\n\tUsage:\n\t\t/pivotal {start,unstart,finish,deliver,accept,reject,unschedule} {story_id}\n\t\t/pivotal query {search_query}\n\tExpected Response: List of stories"
}