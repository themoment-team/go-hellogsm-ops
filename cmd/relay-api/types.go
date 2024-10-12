package main

import (
	"fmt"
)

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

func (e Env) isValid() bool {
	return e == EnvDev || e == EnvProd
}

func (e Env) getEnvName() (Env, error) {
	if e == "" {
		e = EnvDev
	}

	if !e.isValid() {
		return "", fmt.Errorf("잘못된 Env: %s", e)
	}

	return e, nil
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

type Channel string

const (
	Info Channel = "info"
	Mon  Channel = "mon"
)

func (c Channel) isValid() bool {
	return c == Info || c == Mon
}

func (c Channel) getChannelName() (Channel, error) {
	if c == "" {
		c = Info
	}

	if !c.isValid() {
		return "", fmt.Errorf("잘못된 Channel: %s", c)
	}

	return c, nil
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
	Channel     Channel     `json:"channel"`
	Env         Env         `json:"env"`
}
