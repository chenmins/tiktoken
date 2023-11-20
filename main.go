package main

import (
	"encoding/json"
	"fmt"
	"github.com/pkoukk/tiktoken-go"
	tiktokenloader "github.com/pkoukk/tiktoken-go-loader"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// OpenAIResponse 定义接收数据的结构
type OpenAIResponse struct {
	Model   string `json:"model"`
	Object  string `json:"object"`
	Choices []struct {
		Index int `json:"index"`
		Delta struct {
			Content      string `json:"content"`
			FinishReason string `json:"finish_reason"`
			Role         string `json:"role"`
		} `json:"delta"`
	} `json:"choices"`
}

// Statistics 定义统计数据的结构
type Statistics struct {
	ModelName     string
	TotalChars    int
	FinishReasons map[string]int
}

// 流式处理并统计数据的函数
func streamAndCollectStats(dataStream chan string) Statistics {
	stats := Statistics{
		FinishReasons: make(map[string]int),
	}

	// 初始化 tokenizer（根据实际情况调整）
	tiktoken.SetBpeLoader(tiktokenloader.NewOfflineLoader())
	encoding := "gpt-4"

	for data := range dataStream {
		var response OpenAIResponse
		err := json.Unmarshal([]byte(data), &response)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			continue
		}

		if stats.ModelName == "" {
			stats.ModelName = response.Model
		}
		for _, choice := range response.Choices {
			// 用 tokenizer 编码文本
			tkm, err := tiktoken.EncodingForModel(encoding)
			if err != nil {
				fmt.Println("Error getting encoding:", err)
				continue
			}
			tokens := tkm.Encode(choice.Delta.Content, nil, nil)
			if err != nil {
				fmt.Println("Error encoding text:", err)
				continue
			}
			stats.TotalChars += len(tokens)
			if reason := choice.Delta.FinishReason; reason != "" {
				stats.FinishReasons[reason]++
			}
		}
	}

	return stats
}

func main() {
	// 设置 HTTP 路由处理函数
	http.HandleFunc("/encode", encodeHandler)
	http.HandleFunc("/streamTokens", streamTokensHandler)
	// 启动 HTTP 服务器
	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func streamTokensHandler(w http.ResponseWriter, r *http.Request) {
	// 确保是 POST 请求
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// 读取整个请求体
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 按行拆分请求体
	lines := strings.Split(string(body), "\n")

	// 创建一个通道并处理每行数据
	streamChannel := make(chan string)
	go func() {
		for _, line := range lines {
			if strings.HasPrefix(line, "data: ") {
				jsonData := strings.TrimPrefix(line, "data: ")
				streamChannel <- jsonData
			}
		}
		close(streamChannel)
	}()

	// 统计数据
	stats := streamAndCollectStats(streamChannel)
	// 返回统计结果
	json.NewEncoder(w).Encode(stats)
}

func encodeHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应类型为 JSON
	w.Header().Set("Content-Type", "application/json")

	// 读取请求体中的文本
	var data struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 原有的代码逻辑
	tiktoken.SetBpeLoader(tiktokenloader.NewOfflineLoader())
	encoding := "gpt-4"

	tkm, err := tiktoken.EncodingForModel(encoding)
	if err != nil {
		http.Error(w, fmt.Sprintf("getEncoding: %v", err), http.StatusInternalServerError)
		return
	}

	// encode
	token := tkm.Encode(data.Text, nil, nil)

	// 创建响应
	response := struct {
		Tokens    []int `json:"tokens"`
		NumTokens int   `json:"num_tokens"`
	}{
		Tokens:    token,
		NumTokens: len(token),
	}

	// 发送 JSON 响应
	json.NewEncoder(w).Encode(response)
}
