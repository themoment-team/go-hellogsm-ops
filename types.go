package main

type NoticeLevel string

const (
	NoticeLevelError NoticeLevel = "error"
	NoticeLevelWarn  NoticeLevel = "warn"
	NoticeLevelInfo  NoticeLevel = "info"
)

func (s NoticeLevel) getColorCode() int {
	switch s {
	case NoticeLevelError:
		return 0xFF4C4C
	case NoticeLevelWarn:
		return 0xFFA500
	case NoticeLevelInfo:
		return 0x4CAF50
	default:
		return -1
	}
}

func (s NoticeLevel) isValidNoticeLevel() bool {
	switch s {
	case NoticeLevelError, NoticeLevelWarn, NoticeLevelInfo:
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
	Embeds []Embed `json:"embeds"`
}

type Notification struct {
	Title       string      `json:"title"`
	Content     string      `json:"content"`
	NoticeLevel NoticeLevel `json:"noticeLevel"`
}
