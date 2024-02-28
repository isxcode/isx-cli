/*
Copyright © 2023 EchoJamie HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/isxcode/isx-cli/git"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

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

	if projectName == "" {
		fmt.Println("请先选择项目开发")
		os.Exit(1)
	}

	projectPath := viper.GetString(projectName+".dir") + "/" + projectName
	branchName := git.GetCurrentBranchName(projectName, false)

	fmt.Printf(git.BranchTemplate, projectName, branchName, projectPath)
}
