package main

import (
	"encoding/json"
	"fmt"
	"github.com/pkoukk/tiktoken-go"
	tiktokenloader "github.com/pkoukk/tiktoken-go-loader"
	"log"
	"net/http"
)

func main() {
	// 设置 HTTP 路由处理函数
	http.HandleFunc("/encode", encodeHandler)

	// 启动 HTTP 服务器
	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
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
