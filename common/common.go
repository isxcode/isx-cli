/*
Copyright © 2024 jamie HERE <EMAIL ADDRESS>
*/
package common

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"io"
	"os"
)

func HomeDir() string {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return home
}

func CurrentWorkDir() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return dir
}

func Parse(reader io.Reader, v any) {
	body, err := io.ReadAll(reader)
	if err != nil {
		fmt.Println("读取响应体失败:", err)
		os.Exit(1)
	}
	err = json.Unmarshal(body, v)
	if err != nil {
		fmt.Println("解析 JSON 失败:", err)
		os.Exit(1)
	}
}

func ToJsonString(v any) string {
	return string(ToJsonBytes(v))
}

func ToJsonBytes(v any) []byte {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		fmt.Println("解析 JSON 失败:", err)
		os.Exit(1)
	}
	return jsonBytes
}
