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
	// 首先尝试新配置格式 (now-project)
	projectName := viper.GetString("now-project")

	// 如果新配置为空，尝试旧配置格式 (current-project.name)
	if projectName == "" {
		projectName = viper.GetString("current-project.name")
	}

	if projectName == "" {
		fmt.Println("请先使用【isx choose】选择项目")
		os.Exit(1)
	}

	// 获取项目路径 - 支持新旧两种配置格式
	var projectPath string

	// 尝试新配置格式：从 project-list 数组中查找
	type ProjectConfig struct {
		Name string `mapstructure:"name"`
		Dir  string `mapstructure:"dir"`
	}

	var projectList []ProjectConfig
	err := viper.UnmarshalKey("project-list", &projectList)
	if err == nil {
		// 新配置格式：在 project-list 数组中查找项目
		for _, proj := range projectList {
			if proj.Name == projectName {
				projectPath = proj.Dir
				break
			}
		}
	}

	// 如果新配置格式没找到，尝试旧配置格式
	if projectPath == "" {
		projectPath = viper.GetString(projectName + ".dir")
		if projectPath != "" {
			projectPath = projectPath + "/" + projectName
		}
	}

	if projectPath == "" {
		fmt.Printf("项目 %s 未下载，请先使用【isx clone】下载项目代码\n", projectName)
		os.Exit(1)
	}

	branchName := git.GetCurrentBranchName(projectName, projectPath, false)

	fmt.Printf(git.BranchTemplate, projectName, branchName, projectPath)
}
