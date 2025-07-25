package cmd

import (
	"fmt"
	"github.com/isxcode/isx-cli/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"strings"
)

func init() {
	rootCmd.AddCommand(cleanCmd)
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: printCommand("isx clean", 40) + "| 清除项目缓存",
	Long:  `isx clean`,
	Run: func(cmd *cobra.Command, args []string) {
		cleanCmdMain()
	},
}

func cleanCmdMain() {
	// 获取当前项目名称 - 支持新旧配置格式
	projectName := viper.GetString("now-project")
	if projectName == "" {
		projectName = viper.GetString("current-project.name")
	}

	if projectName == "" {
		fmt.Println("请先使用【isx choose】选择项目")
		os.Exit(1)
	}
	var resourcePath string
	if projectName == "spark-yun" {
		resourcePath = "~/.zhiqingyun"
	} else if projectName == "torch-yun" {
		resourcePath = "~/.zhishuyun"
	} else {
		fmt.Println("该项目" + projectName + "暂不支持,请升级isx命令")
		os.Exit(1)
	}

	// 删除前二次确认
	fmt.Printf("是否确认删除该路径(%s)下项目缓存? (y/n) ", common.WhiteText(resourcePath))
	var flag = ""
	fmt.Scanln(&flag)
	flag = strings.Trim(flag, " ")
	flag = strings.ToUpper(flag)
	if flag != "Y" && flag != "N" {
		fmt.Println("输入值异常")
		os.Exit(1)
	}
	if flag == "N" {
		fmt.Println("删除项目缓存已终止")
		os.Exit(0)
	}

	// 删除项目文件
	removeCommand := "rm -rf " + resourcePath
	removeCmd := exec.Command("bash", "-c", removeCommand)
	removeCmd.Stdout = os.Stdout
	removeCmd.Stderr = os.Stderr
	err := removeCmd.Run()
	if err != nil {
		fmt.Println("执行失败:", err)
		os.Exit(1)
	} else {
		fmt.Println(resourcePath + "路径已删除")
	}
}
