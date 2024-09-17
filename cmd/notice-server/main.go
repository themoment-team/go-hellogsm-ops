package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var discordWebhookURL string

func main() {
	initApplicationProperties()

	http.HandleFunc("/notice", handleDiscordWebhook)
	http.HandleFunc("/ping", handlePing)

	log.Println("Server is starting...")

	srv := &http.Server{
		Addr:    ":8085",
		Handler: logRequest(http.DefaultServeMux),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf(err.Error())
		}
	}()

	gracefulShutdown(srv)

	log.Println("Server exiting")
}

// OS에서 command+c 같은 종료 이벤트를 받았을때 server 를 shutdown 하도록 한다.
// graceful shutdown 기능이 있으며, 기존에 listening 중이였던 tcp port 가 kill 된다.
func gracefulShutdown(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server is shutting down...")

	// 기존에 처리되고 있던 요청이 다 처리될때까지 기다린다.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %s", err)
	}
}

func initApplicationProperties() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	discordWebhookURL = os.Getenv("DISCORD_WEBHOOK_URL")
}

// 더모먼트팀 discord 채널로 메시지를 서빙한다.
func handleDiscordWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "지원되지 않는 메서드", http.StatusMethodNotAllowed)
		return
	}

	var notification HellogsmNotification
	if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
		http.Error(w, "JSON 파싱 실패", http.StatusBadRequest)
		return
	}

	if !notification.NoticeLevel.isValidNoticeLevel() {
		http.Error(w, fmt.Sprintf("잘못된 NoticeLevel: %s", notification.NoticeLevel), http.StatusBadRequest)
		return
	}

	if err := sendNotificationToDiscord(notification); err != nil {
		log.Println("디스코드 웹훅 전송 실패", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// ALB 등에서 health check 를 위한 endpoint 를 만든다.
func handlePing(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "지원하지 않는 메소드", http.StatusMethodNotAllowed)
	}
	writer.WriteHeader(http.StatusOK)
}

func sendNotificationToDiscord(notification HellogsmNotification) error {
	embed := Embed{
		Title:       notification.Title,
		Description: notification.Content,
		Color:       notification.NoticeLevel.getColorCode(),
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

// handleDiscordWebhook, handlePing 모두 client 에서 http status code 만으로 성공, 실패 여부를 확인 가능하기에, request 만 로깅한다.
func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
