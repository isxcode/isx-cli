/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
)

var (
	account string
	token   string
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: printCommand("isx login") + "| 登录github账号",
	Long:  `github用户登录`,
	Run: func(cmd *cobra.Command, args []string) {

		// 输入github账号
		fmt.Print("请输入github账号:")
		fmt.Scanln(&account)

		// 输入github令牌
		fmt.Println("快捷链接：https://github.com/settings/tokens")
		fmt.Print("请输入token:")
		fmt.Scanln(&token)

		// 检查令牌是否可用
		checkGithubToken()

		// 保存配置
		saveConfigLogin()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func saveConfigLogin() {
	viper.Set("user.account", account)
	viper.Set("user.token", token)
	viper.WriteConfig()
}

func checkGithubToken() {

	headers := http.Header{}
	headers.Set("Accept", "application/vnd.github+json")
	headers.Set("Authorization", "Bearer "+token)
	headers.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/octocat", nil)

	req.Header = headers
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败:", err)
		os.Exit(1)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// 解析结果
	if resp.StatusCode == http.StatusOK {
		fmt.Println("登录成功，欢迎使用isx-cli开发工具")
	} else {
		fmt.Println("无法验证token合法性，登录失败")
		os.Exit(0)
	}
}
