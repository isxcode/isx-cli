package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/exec"
)

func init() {
	rootCmd.AddCommand(pullCmd)
}

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: printCommand("isx pull <issue_number>", 65) + "| 同步组织代码",
	Long:  `isx pull 123`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			fmt.Println("使用方式不对，请重新输入命令")
			os.Exit(1)
		}

		pullCmdMain(args[0])
	},
}

func pullCmdMain(issueNumber string) {

	branchName := "GH-" + issueNumber

	projectName := viper.GetString("current-project.name")
	projectPath := viper.GetString(projectName+".dir") + "/" + projectName
	rebaseBranch(projectPath, branchName)

	var subRepository []Repository
	viper.UnmarshalKey(viper.GetString("current-project.name")+".sub-repository", &subRepository)
	for _, repository := range subRepository {
		rebaseBranch(projectPath+"/"+repository.Name, branchName)
	}
}

func rebaseBranch(path string, branchName string) {

	// rebase远程的代码
	rebaseCommand := "git fetch upstream && git rebase upstream/" + branchName
	rebaseCmd := exec.Command("bash", "-c", rebaseCommand)
	rebaseCmd.Stdout = os.Stdout
	rebaseCmd.Stderr = os.Stderr
	rebaseCmd.Dir = path
	rebaseCmd.Run()
}
