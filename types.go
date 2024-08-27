package main

type Status string

const (
	StatusError Status = "error"
	StatusWarn  Status = "warn"
	StatusInfo  Status = "info"
)

func (s Status) getColorCode() int {
	switch s {
	case StatusError:
		return 0xFF4C4C
	case StatusWarn:
		return 0xFFA500
	case StatusInfo:
		return 0x4CAF50
	}

	return -1
}

func (s Status) validStatus() bool {
	switch s {
	case StatusError, StatusWarn, StatusInfo:
		return true
	}
	return false
}

type Embed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       int    `json:"color,omitempty"`
}

type DiscordWebhookPayload struct {
	Embeds  []Embed `json:"embeds"`
	Content string  `json:"content"`
}

type Notification struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  Status `json:"status"`
}
