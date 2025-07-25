package cmd

import (
	"fmt"
	"github.com/isxcode/isx-cli/git"
	"github.com/isxcode/isx-cli/github"
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
	Short: printCommand("isx push", 40) + "| 提交代码",
	Long:  `isx push = isx format + isx push`,
	Run: func(cmd *cobra.Command, args []string) {
		pushCmdMain()
	},
}

func pushCmdMain() {
	// 获取当前项目名称 - 支持新旧配置格式
	projectName := viper.GetString("now-project")
	if projectName == "" {
		projectName = viper.GetString("current-project.name")
	}

	if projectName == "" {
		fmt.Println("请先使用【isx choose】选择项目")
		os.Exit(1)
	}

	// 获取项目路径 - 支持新旧配置格式
	var projectPath string

	// 尝试新配置格式：从 project-list 数组中查找
	type ProjectConfig struct {
		Name string `mapstructure:"name"`
		Dir  string `mapstructure:"dir"`
	}

	var projectList []ProjectConfig
	configErr := viper.UnmarshalKey("project-list", &projectList)
	if configErr == nil {
		for _, proj := range projectList {
			if proj.Name == projectName {
				projectPath = proj.Dir
				break
			}
		}
	}

	// 如果新配置格式没找到，尝试旧配置格式
	if projectPath == "" {
		projectDir := viper.GetString(projectName + ".dir")
		if projectDir != "" {
			projectPath = projectDir + "/" + projectName
		}
	}

	if projectPath == "" {
		fmt.Printf("项目 %s 未下载，请先使用【isx clone】下载项目代码\n", projectName)
		os.Exit(1)
	}

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
	branchName := git.GetCurrentBranchName(projectName, projectPath, true)

	// 自动commit 和 提交代码
	commitAndPushCode(projectPath, branchName)

	// 遍历子模块
	subRepository := GetSubRepositories(projectName)
	for _, repository := range subRepository {
		if github.IsRepoForked(viper.GetString("user.account"), repository.Name) {
			commitAndPushCode(projectPath+"/"+repository.Name, branchName)
		}
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
