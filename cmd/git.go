package cmd

import (
	"fmt"
	"github.com/isxcode/isx-cli/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/exec"
)

var gitProjectPath string
var gitProjectName string

func init() {
	rootCmd.AddCommand(gitCmd)
}

var gitCmd = &cobra.Command{
	Use:   "git",
	Short: printCommand("isx git <git_command>", 40) + "| 执行git命令",
	Long:  `isx git <git_command>`,
	Run: func(cmd *cobra.Command, args []string) {
		gitCmdMain(args)
	}, DisableFlagParsing: true,
}

func gitCmdMain(args []string) {
	// 如果没有提供参数，显示 git 状态信息
	if len(args) == 0 {
		fmt.Println("使用方式不对，请输入备注信息：isx git <git_command>")
		os.Exit(1)
	}

	// 获取当前项目名称 - 支持新旧配置格式
	gitProjectName = viper.GetString("now-project")
	if gitProjectName == "" {
		gitProjectName = viper.GetString("current-project.name")
	}

	if gitProjectName == "" {
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
			if proj.Name == gitProjectName {
				projectPath = proj.Dir
				break
			}
		}
	}

	// 如果新配置格式没找到，尝试旧配置格式
	if projectPath == "" {
		gitProjectPath = viper.GetString(gitProjectName + ".dir")
		if gitProjectPath != "" {
			projectPath = gitProjectPath + "/" + gitProjectName
		}
	} else {
		gitProjectPath = projectPath
	}

	if projectPath == "" {
		fmt.Printf("项目 %s 未下载，请先使用【isx clone】下载项目代码\n", gitProjectName)
		os.Exit(1)
	}

	// 进入主项目执行git命令
	gitCmd := exec.Command("git", args...)
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr
	gitCmd.Dir = projectPath
	err := gitCmd.Run()
	if err != nil {
		fmt.Println("执行失败:", err)
		os.Exit(1)
	} else {
		fmt.Println(gitProjectName + " Git命令执行成功")
	}

	// 进入子项目执行命令
	subRepository := GetSubRepositories(gitProjectName)
	for _, repository := range subRepository {

		if github.IsRepoForked(viper.GetString("user.account"), repository.Name) {
			gitCmd := exec.Command("git", args...)
			gitCmd.Stdout = os.Stdout
			gitCmd.Stderr = os.Stderr
			gitCmd.Dir = projectPath + "/" + repository.Name
			err := gitCmd.Run()
			if err != nil {
				fmt.Println("执行失败:", err)
				os.Exit(1)
			} else {
				fmt.Println(repository.Name + " Git命令执行成功")
			}
		}
	}
}
