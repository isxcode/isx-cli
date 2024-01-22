/*
Copyright © 2023 EchoJamie HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os/exec"
	"strings"
)

const nowTmpl = `当前项目: %s
当前分支: %s
项目路径: %s
`

func init() {
	rootCmd.AddCommand(nowCmd)
}

var nowCmd = &cobra.Command{
	Use:   "now",
	Short: printCommand("isx now", 65) + "| 查看当前开发项目",
	Long:  `isx now`,
	Run: func(cmd *cobra.Command, args []string) {
		nowCmdMain()
	},
}

func nowCmdMain() {
	projectName := viper.GetString("current-project.name")
	projectPath := viper.GetString(projectName+".dir") + "/" + projectName

	executeCommand := "git branch --show-current"
	branchCmd := exec.Command("bash", "-c", executeCommand)
	branchCmd.Dir = projectPath
	output, err := branchCmd.Output()
	if err != nil {
		fmt.Printf(nowTmpl, projectName, "获取分支名称失败", projectPath)
		fmt.Println("执行命令失败:", err)
		return
	}
	branchName := strings.Split(string(output), "\n")

	fmt.Printf(nowTmpl, projectName, branchName, projectPath)
}
