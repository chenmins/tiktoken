package main

import (
	"fmt"
	"github.com/pkoukk/tiktoken-go"
	tiktokenloader "github.com/pkoukk/tiktoken-go-loader"
)

func main() {
	// 如果你不想在运行时下载字典，你可以使用离线加载器
	tiktoken.SetBpeLoader(tiktokenloader.NewOfflineLoader())
	text := "Hello, world!"
	//encoding := "gpt-3.5-turbo"
	encoding := "gpt-4"

	tkm, err := tiktoken.EncodingForModel(encoding)
	if err != nil {
		err = fmt.Errorf("getEncoding: %v", err)
		return
	}

	// encode
	token := tkm.Encode(text, nil, nil)

	// tokens
	fmt.Println(token)
	// num_tokens
	fmt.Println(len(token))
}

func main2() {
	text := "Hello, world!"
	encoding := "cl100k_base"

	// 如果你不想在运行时下载字典，你可以使用离线加载器
	tiktoken.SetBpeLoader(tiktokenloader.NewOfflineLoader())
	tke, err := tiktoken.GetEncoding(encoding)
	if err != nil {
		err = fmt.Errorf("getEncoding: %v", err)
		return
	}

	// encode
	token := tke.Encode(text, nil, nil)

	//tokens
	fmt.Println((token))
	// num_tokens
	fmt.Println(len(token))
}
