package robots

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/trinchan/slackbot/robots"
)

type bot struct {
	ProjectID string
	Token     string
}

type searchResultContainer struct {
	Epics   epicsSearchResult   `json:"epics"`
	Query   string              `json:"query"`
	Stories storiesSearchResult `json:"stories"`
}

type epicsSearchResult struct {
	TotalHits         int    `json:"total_hits"`
	Epics             []epic `json:"epics"`
	TotalHitsWithDone int    `json:"total_hits_with_done,omitempty"`
}

type storiesSearchResult struct {
	TotalHits            int     `json:"total_hits"`
	TotalHitsWithDone    int     `json:"total_hits_with_done,omitempty"`
	Stories              []story `json:"stories"`
	TotalPoints          int     `json:"total_points"`
	TotalPointsCompleted int     `json:"total_points_completed,omitempty"`
}

type label struct {
	Kind      string `json:"kind"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Name      string `json:"name"`
	ProjectID int    `json:"project_id"`
	ID        int    `json:"id"`
}

type epic struct {
	CommentIDs  []int  `json:"comment_ids"`
	Name        string `json:"name"`
	Kind        string `json:"kind"`
	ID          int    `json:"id"`
	BeforeID    int    `json:"before_id"`
	UpdatedAt   string `json:"updated_at"`
	URL         string `json:"url"`
	ProjectID   int    `json:"project_id"`
	LabelID     int    `json:"label_id"`
	CreatedAt   string `json:"created_at"`
	Description string `json:"description"`
	AfterID     int    `json:"after_id"`
}

type story struct {
	CommentIDs             []int   `json:"comment_ids,omitempty"`
	CurrentState           string  `json:"current_state"`
	Deadline               string  `json:"deadline,omitempty"`
	RequestedByID          int     `json:"requested_by_id,omitempty"`
	IntegrationID          int     `json:"integration_id,omitempty"`
	Name                   string  `json:"name"`
	OwnedByID              int     `json:"owned_by_id,omitempty"`
	Kind                   string  `json:"kind"`
	Labels                 []label `json:"labels,omitempty"`
	ID                     int     `json:"id"`
	PlannedIterationNumber int     `json:"planned_iteration_number,omitempty"`
	ExternalID             string  `json:"external_id,omitempty"`
	Estimate               int     `json:"estimate"`
	TaskIDs                []int   `json:"task_ids,omitempty"`
	UpdatedAt              string  `json:"updated_at"`
	URL                    string  `json:"url"`
	ProjectID              int     `json:"project_id"`
	StoryType              string  `json:"story_type"`
	AcceptedAt             string  `json:"accepted_at"`
	FollowerIDs            []int   `json:"follower_ids,omitempty"`
	CreatedAt              string  `json:"created_at"`
	Description            string  `json:"description"`
	OwnerIDs               []int   `json:"owner_ids,omitempty"`
}

type task struct {
	Complete    bool   `json:"complete"`
	Kind        string `json:"kind"`
	ID          int    `json:"id"`
	Position    int    `json:"position"`
	UpdatedAt   string `json:"updated_at"`
	StoryID     int    `json:"story_id"`
	CreatedAt   string `json:"created_at"`
	Description string `json:"description"`
}

// Loads the config file and registers the bot with the server for command /${1/(.+)/\L\1/g}.
func init() {
	p := &bot{}
	p.ProjectID = os.Getenv("PIVOTAL_PROJECT_ID")
	p.Token = os.Getenv("PIVOTAL_TOKEN")
	if p.ProjectID == "" || p.Token == "" {
		log.Println("PIVOTAL_PROJECT_ID or PIVOTAL_TOKEN not set, Pivotal bot disabled")
		return
	}
	robots.RegisterRobot("pivotal", p)
}

// All Robots must implement a Run command to be executed when the registered command is received.
func (r bot) Run(p *robots.Payload) (s string) {
	// The string returned here will be shown only to the user who executed the command
	// and will show up as a message from slackbot.
	text := strings.TrimSpace(p.Text)
	if text != "" {
		split := strings.Split(text, " ")
		pc := split[0]
		query := strings.Join(split[1:], " ")
		switch pc {
		case "query":
			return r.query(query)
		case "start", "unstart", "finish", "accept", "reject", "deliver":
			return r.changeState(pc, query)
		}
		return fmt.Sprintf("Unknown pivotal command: %s\n%s", pc, r.Description())
	}
	return ""
}

func (r bot) changeState(state string, storyID string) (result string) {
	params := url.Values{}
	params.Set("current_state", state+"ed")
	req, err := http.NewRequest("PUT", fmt.Sprintf("https://www.pivotaltracker.com/services/v5/projects/%s/stories/%s", r.ProjectID, storyID), nil)
	if err != nil {
		return fmt.Sprintf("ERROR: Error forming put request to Pivotal: %s", err)
	}
	req.URL.RawQuery = params.Encode()
	req.Header.Add("X-TrackerToken", r.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Sprintf("ERROR: Error making put request to Pivotal: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		message := fmt.Sprintf("ERROR: Non-200 Response from Pivotal API: %s", resp.Status)
		log.Println(message)
		return message
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("ERROR: Error reading response body from Pivotal: %s", err)
	}
	s := story{}
	err = json.Unmarshal(contents, &s)
	if err != nil {
		return fmt.Sprintf("ERROR: Couldn't unmarshal pivotal story response into struct: %s", err)
	}
	return fmt.Sprintf("[%s <%s|#%d>] - %s", s.CurrentState, s.URL, s.ID, s.Name)
}

func (r bot) query(query string) (result string) {
	params := url.Values{}
	params.Set("query", query)
	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.pivotaltracker.com/services/v5/projects/%s/search", r.ProjectID), nil)
	req.URL.RawQuery = params.Encode()
	if err != nil {
		return fmt.Sprintf("ERROR: Error forming get request to Pivotal: %s", err)
	}
	req.Header.Add("X-TrackerToken", r.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Sprintf("ERROR: Error making get request to Pivotal: %s", err)
	}
	if resp.StatusCode != 200 {
		message := fmt.Sprintf("ERROR: Non-200 Response from Pivotal API: %s", resp.Status)
		log.Println(message)
		return message
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("ERROR: Error reading response body from Pivotal: %s", err)
	}
	searchResults := searchResultContainer{}
	err = json.Unmarshal(contents, &searchResults)
	if err != nil {
		return fmt.Sprintf("ERROR: Couldn't unmarshal pivotal response into struct: %s", err)
	}
	output := ""
	if searchResults.Stories.TotalHits > 0 {
		for _, story := range searchResults.Stories.Stories {
			output = output + (fmt.Sprintf("[%s <https://www.pivotaltracker.com/s/projects/%d/stories/%d|#%d>] - %s", story.CurrentState, story.ProjectID, story.ID, story.ID, story.Name)) + "\n"
		}
	}
	if output == "" {
		return "No stories matching that query :("
	}
	return strings.TrimSpace(output)
}

func (r bot) Description() (description string) {
	// In addition to a Run method, each Robot must implement a Description method which
	// is just a simple string describing what the Robot does. This is used in the included
	// /c command which gives users a list of commands and descriptions
	return "Interact with the Pivotal API!\n\tUsage:\n\t\t/pivotal {start,unstart,finish,deliver,accept,reject} {story_id}\n\t\t/pivotal query {search_query}\n\tExpected Response: List of stories"
}
