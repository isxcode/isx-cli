/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: printCommand("isx logout", 40) + "| 退出登录",
	Long:  `退出当前登录的GitHub账号`,
	Run: func(cmd *cobra.Command, args []string) {
		logoutCmdMain()
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}

func logoutCmdMain() {
	// 检查是否已经登录
	currentAccount := viper.GetString("user.account")
	currentToken := viper.GetString("user.token")

	if currentAccount == "" && currentToken == "" {
		fmt.Println("当前未登录任何账号")
		return
	}

	// 显示当前登录的账号信息
	if currentAccount != "" {
		fmt.Printf("当前登录账号: %s\n", currentAccount)
	}

	// 确认是否要登出
	fmt.Print("确认要登出吗？(y/n): ")
	var confirm string
	fmt.Scanln(&confirm)

	// 转换为小写进行比较，支持大小写不敏感
	confirm = strings.ToLower(strings.Trim(confirm, " "))

	if confirm != "y" && confirm != "yes" {
		fmt.Println("取消登出")
		return
	}

	// 清除用户登录信息
	clearUserLoginInfo()

	fmt.Println("登出成功！")
}

// clearUserLoginInfo 清除用户登录信息
func clearUserLoginInfo() {
	// 清除账号和token，但保留secret（加密密钥）
	// secret可能在其他地方有用，而且重新生成会导致已加密的数据无法解密
	viper.Set("user.account", "")
	viper.Set("user.token", "")

	// 保存配置
	err := viper.WriteConfig()
	if err != nil {
		fmt.Printf("保存配置失败: %v\n", err)
		return
	}
}
