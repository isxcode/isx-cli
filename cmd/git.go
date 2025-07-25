package cmd

import (
	"fmt"
	"github.com/isxcode/isx-cli/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"strings"
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
		showGitStatus()
		return
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

// showGitStatus 显示 git 状态信息和常用命令
func showGitStatus() {
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

	fmt.Printf("isx git - Git 命令工具 (当前项目: %s)\n", gitProjectName)
	fmt.Println("")

	// 显示当前项目的 git 状态
	fmt.Printf("=== %s 项目状态 ===\n", gitProjectName)
	gitStatusCmd := exec.Command("git", "status", "--short")
	gitStatusCmd.Dir = projectPath
	output, err := gitStatusCmd.Output()
	if err != nil {
		fmt.Printf("获取 git 状态失败: %v\n", err)
	} else {
		if len(output) == 0 {
			fmt.Println("工作目录干净，没有未提交的更改")
		} else {
			fmt.Println("未提交的更改:")
			fmt.Print(string(output))
		}
	}

	// 显示当前分支
	gitBranchCmd := exec.Command("git", "branch", "--show-current")
	gitBranchCmd.Dir = projectPath
	branchOutput, err := gitBranchCmd.Output()
	if err != nil {
		fmt.Printf("获取当前分支失败: %v\n", err)
	} else {
		fmt.Printf("当前分支: %s\n", strings.TrimSpace(string(branchOutput)))
	}

	fmt.Println("")
	fmt.Println("使用方法:")
	fmt.Println("  isx git <git_command>")
	fmt.Println("")
	fmt.Println("常用命令:")
	fmt.Println("  isx git status     - 查看状态")
	fmt.Println("  isx git add .      - 添加所有更改")
	fmt.Println("  isx git commit -m \"message\" - 提交更改")
	fmt.Println("  isx git push       - 推送到远程")
	fmt.Println("  isx git pull       - 拉取远程更改")
	fmt.Println("  isx git branch     - 查看分支")
	fmt.Println("  isx git log --oneline - 查看提交历史")
}
