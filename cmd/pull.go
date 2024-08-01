package cmd

import (
	"github.com/isxcode/isx-cli/git"
	"github.com/isxcode/isx-cli/github"
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
	Short: printCommand("isx pull", 65) + "| 拉取组织代码",
	Long:  `isx pull`,
	Run: func(cmd *cobra.Command, args []string) {

		pullCmdMain()
	},
}

func pullCmdMain() {

	branchName := git.GetCurrentBranchName(viper.GetString("current-project.name"), true)

	projectName := viper.GetString("current-project.name")
	projectPath := viper.GetString(projectName+".dir") + "/" + projectName
	rebaseBranch(projectPath, branchName)

	var subRepository []Repository
	viper.UnmarshalKey(viper.GetString("current-project.name")+".sub-repository", &subRepository)
	for _, repository := range subRepository {
		if github.IsRepoForked(viper.GetString("user.account"), repository.Name) {
			rebaseBranch(projectPath+"/"+repository.Name, branchName)
		}
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
