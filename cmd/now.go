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
	Short: printCommand("isx now", 65) + "| 查看项目信息",
	Long:  `isx now`,
	Run: func(cmd *cobra.Command, args []string) {
		nowCmdMain()
	},
}

func nowCmdMain() {
	projectName := viper.GetString("current-project.name")

	if projectName == "" {
		fmt.Println("请先使用【isx choose】选择项目")
		os.Exit(1)
	}

	projectPath := viper.GetString(projectName+".dir") + "/" + projectName
	branchName := git.GetCurrentBranchName(projectName, false)

	fmt.Printf(git.BranchTemplate, projectName, branchName, projectPath)
}
