package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/exec"
	"strings"
)

func init() {
	rootCmd.AddCommand(pullCmd)
}

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: printCommand("isx pull <issue_number>", 65) + "| 同步组织代码",
	Long:  `isx pull 123`,
	Run: func(cmd *cobra.Command, args []string) {

		pullCmdMain()
	},
}

func pullCmdMain() {

	branchName := getBranchName()

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

func getBranchName() string {

	projectName := viper.GetString("current-project.name")
	projectPath := viper.GetString(projectName+".dir") + "/" + projectName

	executeCommand := "git branch --show-current"
	branchCmd := exec.Command("bash", "-c", executeCommand)
	branchCmd.Dir = projectPath
	output, err := branchCmd.Output()
	if err != nil {
		fmt.Printf(nowTmpl, projectName, "获取分支名称失败", projectPath)
		fmt.Println("执行命令失败:", err)
		log.Fatal(err)
		os.Exit(1)
	}
	return strings.Split(string(output), "\n")[0]
}
