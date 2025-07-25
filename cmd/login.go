/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/isxcode/isx-cli/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	accountOrEmail string
	token          string
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: printCommand("isx login", 40) + "| 用户登录",
	Long:  `github用户登录，支持用户名和邮箱`,
	Run: func(cmd *cobra.Command, args []string) {
		loginCmdMain()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func loginCmdMain() {
	// 检查是否已经登录
	if isAlreadyLoggedIn() {
		currentAccount := viper.GetString("user.account")
		fmt.Printf("当前已登录账号: %s\n", currentAccount)
		fmt.Println("如需切换账号，请先使用 【isx logout】 退出当前账号")
		return
	}

	// 输入github令牌
	fmt.Println("快捷链接：https://github.com/settings/tokens")
	fmt.Println("请输入GitHub Personal Access Token:")
	fmt.Scanln(&token)

	// 检查令牌是否可用并获取用户信息
	userInfo := checkGithubTokenAndGetUserInfo()

	// 保存配置
	saveConfigLogin(userInfo)
}

// isAlreadyLoggedIn 检查用户是否已经登录
func isAlreadyLoggedIn() bool {
	currentAccount := viper.GetString("user.account")
	currentToken := viper.GetString("user.token")

	// 如果账号和token都不为空，则认为已登录
	return currentAccount != "" && currentToken != ""
}

func saveConfigLogin(userInfo *common.GitHubUserInfo) {
	// 优先使用用户名，如果没有用户名则使用邮箱
	account := userInfo.Login
	if account == "" && userInfo.Email != "" {
		account = userInfo.Email
	}

	// 加密token后保存
	encryptedToken := common.Encrypt(token)

	viper.Set("user.account", account)
	viper.Set("user.token", encryptedToken)
	viper.WriteConfig()

	fmt.Printf("登录成功！\n")
	fmt.Printf("用户名: %s\n", userInfo.Login)
	if userInfo.Name != "" {
		fmt.Printf("姓名: %s\n", userInfo.Name)
	}
	if userInfo.Email != "" {
		fmt.Printf("邮箱: %s\n", userInfo.Email)
	}
	fmt.Println("欢迎使用isx-cli开发工具")
}

func checkGithubTokenAndGetUserInfo() *common.GitHubUserInfo {
	// 首先检查token是否有效
	if !common.CheckUserAccount(token) {
		fmt.Println("无法验证token合法性，登录失败")
		os.Exit(1)
	}

	// 获取用户信息
	userInfo, err := common.GetGitHubUserInfo(token)
	if err != nil {
		fmt.Printf("获取用户信息失败: %v\n", err)
		os.Exit(1)
	}

	return userInfo
}
