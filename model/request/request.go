package request

type SlackRequest struct {
	Channel   string  `json:"channel,omitempty"`
	Type      string  `json:"type,omitempty"`
	Blocks    []Block `json:"blocks,omitempty"`
	IconEmoji string  `json:"icon_emoji,omitempty"`
	Username  string  `json:"username,omitempty"`
}

type Block struct {
	Type   string     `json:"type,omitempty"`
	Text   *Markdown  `json:"text,omitempty"`
	Fields []Markdown `json:"fields,omitempty"`
}

type Markdown struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
}
