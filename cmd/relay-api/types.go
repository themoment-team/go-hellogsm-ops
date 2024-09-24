package main

import "os"

type NoticeLevel string

const (
	NoticeLevelError NoticeLevel = "error"
	NoticeLevelWarn  NoticeLevel = "warn"
	NoticeLevelInfo  NoticeLevel = "info"
)

func (s NoticeLevel) isValidNoticeLevel() bool {
	switch s {
	case NoticeLevelError, NoticeLevelWarn, NoticeLevelInfo:
		return true
	default:
		return false
	}
}

type Env string

const (
	EnvDev  Env = "dev"
	EnvProd Env = "prod"
)

func (e Env) handleEnv() bool {
	if e == "" {
		e = EnvDev
	}

	var envName string
	switch e {
	case EnvDev:
		envName = "DEV_DISCORD_WEBHOOK_URL"
	case EnvProd:
		envName = "PROD_DISCORD_WEBHOOK_URL"
	default:
		return false
	}

	discordWebhookURL = os.Getenv(envName)
	return true
}

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

type Embed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       int    `json:"color,omitempty"`
}

type DiscordWebhookPayload struct {
	Embeds []Embed `json:"embeds"`
}

type HellogsmNotification struct {
	Title       string      `json:"title"`
	Content     string      `json:"content"`
	NoticeLevel NoticeLevel `json:"noticeLevel"`
	Env         Env         `json:"env"`
}
