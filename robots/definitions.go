package robots

type SlashCommand struct {
	Payload
	Command string `schema:"command"`
}

type Payload struct {
	Token       string `schema:"token"`
	TeamID      string `schema:"team_id"`
	ChannelID   string `schema:"channel_id"`
	ChannelName string `schema:"channel_name"`
	UserID      string `schema:"user_id"`
	UserName    string `schema:"user_name"`
	Text        string `schema:"text,omitempty"`
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

type IncomingWebhook struct {
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

type Configuration struct {
	Domain string `schema:"domain"`
	Port   int    `schema:"port"`
	Token  string `schema:"token"`
}

type Robot interface {
	Run(p *Payload) (botString string)
	Description() (description string)
}
