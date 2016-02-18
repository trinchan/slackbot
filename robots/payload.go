package robots

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type SlashCommand struct {
	Payload
	Command string `schema:"command"`
}

type Payload struct {
	Token       string  `schema:"token"`
	TeamID      string  `schema:"team_id"`
	TeamDomain  string  `schema:"team_domain,omitempty"`
	ChannelID   string  `schema:"channel_id"`
	ChannelName string  `schema:"channel_name"`
	Timestamp   float64 `schema:"timestamp,omitempty"`
	UserID      string  `schema:"user_id"`
	UserName    string  `schema:"user_name"`
	Text        string  `schema:"text,omitempty"`
	TriggerWord string  `schema:"trigger_word,omitempty"`
	ServiceID   string  `schema:"service_id,omitempty"`
	ResponseUrl string  `schema:"response_url,omitempty"`
	BotID       string  `schema:"bot_id,omitempty"`
	BotName     string  `schema:"bot_name,omitempty"`
	Robot       string
}

type OutgoingWebHook struct {
	Payload
	TriggerWord string `schema:"trigger_word"`
}

type OutgoingWebHookResponse struct {
	Text      string     `json:"text"`
	Parse     ParseStyle `json:"parse,omitempty"`
	LinkNames bool       `json:"link_names,omitempty"`
	Markdown  bool       `json:"mrkdwn,omitempty"`
}

type ParseStyle string

var (
	ParseStyleFull = ParseStyle("full")
	ParseStyleNone = ParseStyle("none")
)

type Message struct {
	Domain      string       `json:"domain"`
	Channel     string       `json:"channel"`
	Username    string       `json:"username"`
	Text        string       `json:"text"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	IconURL     string       `json:"icon_url,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
	UnfurlLinks bool         `json:"unfurl_links,omitempty"`
	Parse       ParseStyle   `json:"parse,omitempty"`
	LinkNames   bool         `json:"link_names,omitempty"`
	Markdown    bool         `json:"mrkdwn,omitempty"`
}

type IncomingWebhook Message
type SlashCommandResponse Message

type Attachment struct {
	Fallback   string            `json:"fallback"`
	Pretext    string            `json:"pretext,omitempty"`
	Text       string            `json:"text,omitempty"`
	Color      string            `json:"color,omitempty"`
	Fields     []AttachmentField `json:"fields,omitempty"`
	MarkdownIn []MarkdownField   `json:"mrkdown_in,omitempty"`
}

type MarkdownField string

var (
	MarkdownFieldPretext  = MarkdownField("pretext")
	MarkdownFieldText     = MarkdownField("text")
	MarkdownFieldTitle    = MarkdownField("title")
	MarkdownFieldFields   = MarkdownField("fields")
	MarkdownFieldFallback = MarkdownField("fallback")
)

type AttachmentField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short,omitempty"`
}

// Send uses the IncomingWebhook API to post a message to a slack channel
func (i IncomingWebhook) Send() error {
	u := os.Getenv(fmt.Sprintf("%s_IN_URL", strings.ToUpper(i.Domain)))
	if u == "" {
		return fmt.Errorf("Slack Incoming Webhook URL not found for domain %s (check %s)", i.Domain, fmt.Sprintf("%s_IN_URL", strings.ToUpper(i.Domain)))
	}
	return Message(i).sendToUrl(u)
}

// Send a response to the ResponseUrl in the Payload
func (r SlashCommandResponse) Send(p *Payload) error {
	if p.ResponseUrl == "" {
		return fmt.Errorf("Empty ResponseUrl in Payload: %v", p)
	}
	return Message(r).sendToUrl(p.ResponseUrl)
}

func (i Message) sendToUrl(u string) error {
	if u == "" {
		return fmt.Errorf("Empty URL")
	}
	webhook, err := url.Parse(u)
	if err != nil {
		log.Printf("Error parsing URL \"%s\": %v", u, err)
		return err
	}

	p, err := json.Marshal(i)
	if err != nil {
		return err
	}

	data := url.Values{}
	data.Set("payload", string(p))

	webhook.RawQuery = data.Encode()
	resp, err := http.PostForm(webhook.String(), data)
	if resp.StatusCode != 200 {
		message := fmt.Sprintf("ERROR: Non-200 Response from Slack URL \"%s\": %s", u, resp.Status)
		log.Println(message)
	}
	return err
}
