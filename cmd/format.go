package cmd

import (
	"fmt"
	"github.com/isxcode/isx-cli/git"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/exec"
)

func init() {
	rootCmd.AddCommand(formatCmd)
}

var formatCmd = &cobra.Command{
	Use:   "format",
	Short: printCommand("isx format", 65) + "| 代码格式化",
	Long:  `isx format`,
	Run: func(cmd *cobra.Command, args []string) {
		formatCmdMain()
	},
}

func formatCmdMain() {

	projectName := viper.GetString("current-project.name")
	projectPath := viper.GetString(projectName+".dir") + "/" + projectName

	// 除了isx-cli项目，其他都要使用gradle 格式化代码
	if "isx-cli" != projectName {
		gradleCmd := exec.Command("./gradlew", "format")
		gradleCmd.Stdout = os.Stdout
		gradleCmd.Stderr = os.Stderr
		gradleCmd.Dir = projectPath
		err := gradleCmd.Run()
		if err != nil {
			fmt.Println("执行失败:", err)
			os.Exit(1)
		} else {
			fmt.Println("执行成功")
		}
	}

	// 获取当前分支
	branchName := git.GetCurrentBranchName(viper.GetString("current-project.name"), true)

	// 自动commit 和 提交代码
	commitAndPushCode(projectPath, branchName)

	// 遍历子模块
	var subRepository []Repository
	viper.UnmarshalKey(viper.GetString("current-project.name")+".sub-repository", &subRepository)
	for _, repository := range subRepository {
		commitAndPushCode(projectPath+"/"+repository.Name, branchName)
	}

}

func commitAndPushCode(path string, branchName string) {

	gitAddCommand := "git add ."
	gitAddCmd := exec.Command("bash", "-c", gitAddCommand)
	gitAddCmd.Stdout = os.Stdout
	gitAddCmd.Stderr = os.Stderr
	gitAddCmd.Dir = path
	err := gitAddCmd.Run()
	if err != nil {
		fmt.Println("git add . 异常")
	}

	gitCommitCommand := "git commit -m '格式化代码'"
	gitCommitCmd := exec.Command("bash", "-c", gitCommitCommand)
	gitCommitCmd.Stdout = os.Stdout
	gitCommitCmd.Stderr = os.Stderr
	gitCommitCmd.Dir = path
	err = gitCommitCmd.Run()
	if err != nil {
		fmt.Println("无代码需要commit")
	}

	// 推送代码
	pushOriginCommand := "git push origin " + branchName
	pushOriginCmd := exec.Command("bash", "-c", pushOriginCommand)
	pushOriginCmd.Stdout = os.Stdout
	pushOriginCmd.Stderr = os.Stderr
	pushOriginCmd.Dir = path
	err = pushOriginCmd.Run()
	if err != nil {
		fmt.Println("无法推送，请谨慎尝试强推： isx git push origin " + branchName + " -f")
	}
}
