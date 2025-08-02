package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/ledongthuc/pdf"
	"github.com/sashabaranov/go-openai"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const chunkSize = 2000

type requestPayload struct {
	FilePath string `json:"file_path"`
}

func extractTextFromPDF(filePath string) (string, error) {
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var content string
	totalPages := r.NumPage()
	for i := 1; i <= totalPages; i++ {
		p := r.Page(i)
		if p.V.IsNull() {
			continue
		}
		text, _ := p.GetPlainText(nil)
		content += text
	}
	return content, nil
}

func splitText(text string, size int) []string {
	runes := []rune(text)
	var chunks []string
	for start := 0; start < len(runes); start += size {
		end := start + size
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[start:end]))
	}
	return chunks
}

func callChatGPT(client *openai.Client, text string) (string, error) {
	prompt := fmt.Sprintf(`
다음 문서에서 다음 항목들을 정확하게 추출해줘:

- 이름
- 학년
- 학년별 과목 점수 (예: {"국어": "95", "수학": "88", ...})

결과는 아래 조건에 맞는 JSON 형식으로 출력해줘:

- key는 모두 camelCase로
- value는 모두 string 형식으로
- 학년별 과목 점수는 과목명이 key인 내부 객체 형태로

문서 내용:
%s`, text)

	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	}

	var resp openai.ChatCompletionResponse
	var err error
	retries := 3

	for i := 0; i < retries; i++ {
		resp, err = client.CreateChatCompletion(context.Background(), req)
		if err == nil {
			break
		}
		if apiErr, ok := err.(*openai.APIError); ok && apiErr.HTTPStatusCode == 429 {
			wait := time.Duration((i+1)*5) * time.Second
			log.Printf("429 오류. %d초 후 재시도...\n", wait/time.Second)
			time.Sleep(wait)
			continue
		}
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("응답에 내용이 없습니다")
	}
	return resp.Choices[0].Message.Content, nil
}

func processPDFHandler(client *openai.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST 요청만 허용됩니다", http.StatusMethodNotAllowed)
			return
		}

		var req requestPayload
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "잘못된 요청 형식", http.StatusBadRequest)
			return
		}

		if req.FilePath == "" {
			http.Error(w, "file_path가 필요합니다", http.StatusBadRequest)
			return
		}

		text, err := extractTextFromPDF(req.FilePath)
		if err != nil {
			http.Error(w, fmt.Sprintf("PDF 읽기 실패: %v", err), http.StatusInternalServerError)
			return
		}

		chunks := splitText(text, chunkSize)
		var results []map[string]interface{}

		for _, chunk := range chunks {
			rawResult, err := callChatGPT(client, chunk)
			if err != nil {
				log.Printf("ChatGPT 처리 실패: %v", err)
				continue
			}

			// JSON 파싱 처리
			dec := json.NewDecoder(strings.NewReader(rawResult))
			dec.DisallowUnknownFields()

			for dec.More() {
				var obj map[string]interface{}
				if err := dec.Decode(&obj); err != nil {
					log.Printf("JSON 파싱 실패: %v", err)
					continue
				}
				results = append(results, obj)
			}

			time.Sleep(2 * time.Second)
		}

		// 응답
		response := map[string]interface{}{
			"status": "ok",
			"data":   results,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env 파일 로드 실패")
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY 환경변수가 설정되어 있지 않습니다.")
	}

	client := openai.NewClient(apiKey)

	http.HandleFunc("/process-pdf", processPDFHandler(client))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
