/*
Copyright © 2024 jamie HERE <EMAIL ADDRESS>
*/
package github

import (
	"fmt"
	"github.com/isxcode/isx-cli/common"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
)

func Get(url string, reqBody io.Reader) *http.Response {
	return request(url, "GET", nil)
}

func Post(url string, reqBody io.Reader) *http.Response {
	return request(url, "POST", reqBody)
}

func request(url, method string, reqBody io.Reader) *http.Response {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, reqBody)

	req.Header = common.GitHubHeader(viper.GetString("user.token"))
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败:", err)
		os.Exit(1)
	}
	return resp
}

/*
	CloseRespBody 关闭响应体
	@Example
		defer closeRespBody(resp.Body)
*/

func CloseRespBody(Body io.ReadCloser) {
	err := Body.Close()
	if err != nil {
		fmt.Println("关闭响应体失败:", err)
	}
}
