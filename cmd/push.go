package cmd

import (
	"fmt"
	"github.com/isxcode/isx-cli/git"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"runtime"
)

var forceFlag bool

func init() {
	pushCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Force push")
	rootCmd.AddCommand(pushCmd)
}

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: printCommand("isx push", 65) + "| 格式化代码后,提交代码",
	Long:  `isx push = isx format + isx push`,
	Run: func(cmd *cobra.Command, args []string) {
		pushCmdMain()
	},
}

func pushCmdMain() {

	projectName := viper.GetString("current-project.name")
	projectDir := viper.GetString(projectName + ".dir")
	projectPath := projectDir + "/" + projectName

	// 除了isx-cli项目，其他都要使用gradle 格式化代码
	if "isx-cli" != projectName && "tools-yun" != projectName {

		var gradleCmd *exec.Cmd
		if runtime.GOOS == "windows" {
			gradleCmd = exec.Command("bash", "-c", "./gradlew.bat format")
		} else {
			gradleCmd = exec.Command("./gradlew", "format")
		}

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
	pushOriginCommand := ""
	if forceFlag {
		pushOriginCommand = "git push origin " + branchName + " -f"
	} else {
		pushOriginCommand = "git push origin " + branchName
	}
	pushOriginCmd := exec.Command("bash", "-c", pushOriginCommand)
	pushOriginCmd.Stdout = os.Stdout
	pushOriginCmd.Stderr = os.Stderr
	pushOriginCmd.Dir = path
	err = pushOriginCmd.Run()
	if err != nil {
		fmt.Println("无法推送，请谨慎尝试强推： isx git push origin " + branchName + " -f")
	}
}
