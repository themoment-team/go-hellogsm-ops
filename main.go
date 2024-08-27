package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var discordWebhookURL string

func sendToDiscordWebhook(notification Notification) error {
	embed := Embed{
		Title:       notification.Title,
		Description: notification.Content,
		Color:       notification.Status.getColorCode(),
	}

	payload := DiscordWebhookPayload{
		Embeds: []Embed{embed},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("payload 변환 실패: %w", err)
	}

	res, err := http.Post(discordWebhookURL, "application/json", bytes.NewReader(payloadBytes))
	if err != nil {
		return fmt.Errorf("HTTP 요청 실패: %w", err)
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent {
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("디스코드 웹훅 응답 오류: %s. 처리 중 응답 읽기 실패: %w", res.Status, err)
		}

		return fmt.Errorf("디스코드 웹훅 응답 오류: %s - %s", res.Status, body)
	}

	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "지원되지 않는 메서드", http.StatusMethodNotAllowed)
		return
	}

	var notification Notification
	if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
		http.Error(w, "JSON 파싱 실패", http.StatusBadRequest)
		return
	}

	if !notification.Status.validStatus() {
		http.Error(w, fmt.Sprintf("잘못된 Status: %s", notification.Status), http.StatusBadRequest)
		return
	}

	if err := sendToDiscordWebhook(notification); err != nil {
		log.Println("디스코드 웹훅 전송 실패", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	discordWebhookURL = os.Getenv("DISCORD_WEBHOOK_URL")

	http.HandleFunc("/notice", handler)

	if err := http.ListenAndServe(":8085", nil); err != nil {
		panic(err)
	}
}
