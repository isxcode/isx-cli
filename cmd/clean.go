package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/exec"
)

func init() {
	rootCmd.AddCommand(cleanCmd)
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: printCommand("isx clean", 65) + "| 删除项目缓存",
	Long:  `isx clean`,
	Run: func(cmd *cobra.Command, args []string) {
		cleanCmdMain()
	},
}

func cleanCmdMain() {

	projectName := viper.GetString("current-project.name")
	resourcePath := ""
	if projectName == "spark-yun" {
		resourcePath = "zhiqingyun"
	} else if projectName == "flink-yun" {
		resourcePath = "zhiliuyun"
	} else {
		fmt.Println("该项目" + projectName + "暂不支持,请升级isx命令")
		os.Exit(1)
	}

	// 删除项目文件
	removeCommand := "rm -rf ~/." + resourcePath
	removeCmd := exec.Command("bash", "-c", removeCommand)
	removeCmd.Stdout = os.Stdout
	removeCmd.Stderr = os.Stderr
	err := removeCmd.Run()
	if err != nil {
		fmt.Println("执行失败:", err)
		os.Exit(1)
	} else {
		fmt.Println("~/." + resourcePath + "路径已删除")
	}
}
