/*
Copyright © 2024 jamie HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"strings"
)

var ToolboxScriptPath = "/Users/jamie/Library/Application Support/JetBrains/Toolbox/scripts"
var VscodeCmdPath = "/usr/local/bin/code"

var projectTypeMap = map[string]string{
	"isx-cli":   "goland",
	"spark-yun": "idea",
	"tools-yun": "code",
}

func init() {
	rootCmd.AddCommand(openCmd)
}

var openCmd = &cobra.Command{
	Use:   "open",
	Short: printCommand("isx open", 65) + "| 使用IDE打开当前项目",
	Long:  `isx open [project-name]`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Println("使用方式不对，请重新输入命令")
			os.Exit(1)
		}
		openCmdMain(args)
	},
}

func openCmdMain(args []string) {
	if len(args) == 1 {
		openCmdSpecified(args[0])
	}
	openNow()
}

func openNow() {
	projectName := viper.GetString("current-project.name")
	if projectName == "" {
		fmt.Println("请先使用【isx choose】选择项目")
		os.Exit(1)
	}
	openCmdSpecified(projectName)
}

func openCmdSpecified(projectName string) {

	projectDir := viper.GetString(projectName + ".dir")
	if projectDir == "" {
		fmt.Println("请检查项目名称是否正确")
		os.Exit(1)
	}

	projectPath := projectDir + "/" + projectName
	ide := projectTypeMap[projectName]
	if checkIde(ide) && execOpen(ide, projectPath) != nil {

		wholeCmdPath := VscodeCmdPath
		if ide != "code" {
			wholeCmdPath = ToolboxScriptPath + "/" + ide
		}

		err := execOpen(wholeCmdPath, projectPath)
		if err != nil {
			fmt.Println("命令执行失败", err)
			fmt.Println("")
			os.Exit(1)
		}
	}
}

func execOpen(command string, projectPath string) error {
	openCmd := exec.Command(command, projectPath)
	openCmd.Stdout = os.Stdout
	openCmd.Stderr = os.Stderr
	err := openCmd.Run()
	if err == nil {
		fmt.Println("open命令执行成功")
		os.Exit(0)
	}
	return err
}

func checkIde(ideName string) bool {
	cmd := exec.Command("bash", "-c", "echo $PATH")
	output, _ := cmd.Output()
	if ideName == "code" {
		return strings.Contains(string(output), "/usr/local/bin")
	}
	return strings.Contains(string(output), ToolboxScriptPath)
}
