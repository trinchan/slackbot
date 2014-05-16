package robots

type SlashCommand struct {
	Token        string `schema:"token"`
	Team_ID      string `schema:"team_id"`
	Channel_ID   string `schema:"channel_id"`
	Channel_Name string `schema:"channel_name"`
	User_ID      string `schema:"user_id"`
	User_Name    string `schema:"user_name"`
	Command      string `schema:"command"`
	Text         string `schema:"text,omitempty"`
    Trigger_Word string `schema:"trigger_word,omitempty"`
    Team_Domain  string `schema:"team_domain,omitempty"`
    Service_ID   string `schema:"service_id,omitempty"`
    Timestamp    float64 `schema:"timestamp,omitempty"`
}

type IncomingWebhook struct {
	Channel      string       `json:"channel"`
	Username     string       `json:"username"`
	Text         string       `json:"text"`
	Icon_Emoji   string       `json:"icon_emoji,omitempty"`
	Icon_URL     string       `json:"icon_url,omitempty"`
	Attachments  []Attachment `json:"attachments,omitempty"`
	Unfurl_Links bool         `json:"unfurl_links,omitempty"`
	Parse        string       `json:"parse,omitempty"`
	Link_Names   bool         `json:"link_names,omitempty"`
}

type Attachment struct {
	Fallback string            `json:"fallback"`
	Pretext  string            `json:"pretext,omitempty"`
	Text     string            `json:"text,omitempty"`
	Color    string            `json:"color,omitempty"`
	Fields   []AttachmentField `json:fields,omitempty`
}

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
	Run(command *SlashCommand) (slashCommandImmediateReturn string)
	Description() (description string)
}
