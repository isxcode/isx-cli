package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"runtime"
)

func init() {
	rootCmd.AddCommand(formatCmd)
}

var formatCmd = &cobra.Command{
	Use:   "format",
	Short: printCommand("isx format", 65) + "| 格式化代码",
	Long:  `isx format`,
	Run: func(cmd *cobra.Command, args []string) {
		formatCmdMain()
	},
}

func formatCmdMain() {
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
		projectPath = viper.GetString(projectName + ".dir")
		if projectPath != "" {
			projectPath = projectPath + "/" + projectName
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
}
