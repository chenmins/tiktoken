package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/pkoukk/tiktoken-go"
	tiktokenloader "github.com/pkoukk/tiktoken-go-loader"
)

// OpenAIRequest 定义接收请求体的结构
type OpenAIRequest struct {
	Request struct {
		Body string `json:"body"` // Body 是一个 JSON 编码的字符串
	} `json:"request"`
	Response struct {
		Headers map[string]string `json:"headers"`
		Body    string            `json:"body"`
	} `json:"response"`
}

type Message struct {
	Role     string      `json:"role"`
	Content  string      `json:"content"`
	Metadata interface{} `json:"metadata"`
	Tools    interface{} `json:"tools"`
}

type Choice struct {
	Index        int         `json:"index"`
	Message      Message     `json:"message"`
	FinishReason string      `json:"finish_reason"`
	History      interface{} `json:"history"`
}

type TokenResponse struct {
	Model            string   `json:"model"`
	Object           string   `json:"object"`
	PromptTokens     int      `json:"prompt_tokens"`
	CompletionTokens int      `json:"completion_tokens"`
	TotalTokens      int      `json:"total_tokens"`
	Choices          []Choice `json:"choices"`
}

func main() {
	// 初始化 tokenizer
	tiktoken.SetBpeLoader(tiktokenloader.NewOfflineLoader())

	http.HandleFunc("/", handler)
	log.Println("Server starting on port 8888...")
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req OpenAIRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 检查 Content-Type 是否为 text/event-stream
	contentType, ok := req.Response.Headers["content-type"]
	if ok && strings.Contains(contentType, "text/event-stream") {
		// 解码 request.body
		var requestBody struct {
			Messages []Message `json:"messages"`
		}
		err = json.Unmarshal([]byte(req.Request.Body), &requestBody)
		if err != nil {
			http.Error(w, "Failed to decode body: "+err.Error(), http.StatusBadRequest)
			return
		}

		// 拼接 messages 中的 role 和 content 并计算 tokens
		var combinedPrompt string
		for _, msg := range requestBody.Messages {
			combinedPrompt += msg.Role + ": " + msg.Content + " " // 根据需要调整拼接格式
		}
		promptTokens, err := calculateTokens(combinedPrompt)
		if err != nil {
			http.Error(w, "Failed to calculate prompt tokens: "+err.Error(), http.StatusInternalServerError)
			return
		}

		model, completionTokens, combinedChoice, err := calculateCompletionTokens(req.Response.Body)
		if err != nil {
			http.Error(w, "Failed to calculate completion tokens: "+err.Error(), http.StatusInternalServerError)
			return
		}

		resp := TokenResponse{
			Model:            model,
			Object:           "chat.completion",
			PromptTokens:     promptTokens,
			CompletionTokens: completionTokens,
			TotalTokens:      promptTokens + completionTokens,
			Choices:          []Choice{combinedChoice},
		}

		json.NewEncoder(w).Encode(resp)
	} else {
		// 如果不是 text/event-stream，则直接返回 response.body
		w.Write([]byte(req.Response.Body))
	}

}

func calculateTokens(text string) (int, error) {
	encoding := "gpt-3.5-turbo"
	tkm, err := tiktoken.EncodingForModel(encoding)
	if err != nil {
		return 0, err
	}

	tokens := tkm.Encode(text, nil, nil)
	tokenCount := len(tokens)

	// 添加日志记录
	log.Printf("Processed Text: %s\n", text)
	log.Printf("Token Count: %d\n", tokenCount)

	return tokenCount, nil
}

func calculateCompletionTokens(responseBody string) (string, int, Choice, error) {
	lines := strings.Split(responseBody, "\n")
	var combinedContent, model, lastFinishReason, lastRole string

	for _, line := range lines {
		if strings.HasPrefix(line, "data: ") {
			if strings.TrimSpace(line) == "data: [DONE]" {
				continue
			}

			var response struct {
				Model   string `json:"model"`
				Choices []struct {
					Index int `json:"index"`
					Delta struct {
						Content string `json:"content"`
						Role    string `json:"role"`
					} `json:"delta"`
					FinishReason string `json:"finish_reason"`
				} `json:"choices"`
			}

			jsonData := strings.TrimPrefix(line, "data: ")
			err := json.Unmarshal([]byte(jsonData), &response)
			if err != nil {
				log.Printf("Error parsing JSON data: %v", err)
				continue
			}

			if model == "" {
				model = response.Model
			}

			for _, choice := range response.Choices {
				combinedContent += choice.Delta.Content
				// 只在 finish_reason 有非空值时更新
				if choice.FinishReason != "" {
					lastFinishReason = choice.FinishReason
				}
				if choice.Delta.Role != "" {
					lastRole = choice.Delta.Role
				}
			}
		}
	}

	tokenCount, err := calculateTokens(combinedContent)
	combinedChoice := Choice{
		Index: 0,
		Message: Message{
			Role:    lastRole,
			Content: combinedContent,
		},
		FinishReason: lastFinishReason,
	}

	return model, tokenCount, combinedChoice, err
}
