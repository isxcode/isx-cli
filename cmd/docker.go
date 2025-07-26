package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/exec"
)

func init() {
	rootCmd.AddCommand(dockerCmd)
}

var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: printCommand("isx docker", 40) + "| 构建Docker镜像",
	Long:  `构建项目的Docker镜像`,
	Run: func(cmd *cobra.Command, args []string) {
		dockerCmdMain()
	},
}

func dockerCmdMain() {
	// 获取当前项目名称 - 支持新旧配置格式
	projectName := viper.GetString("now-project")
	if projectName == "" {
		projectName = viper.GetString("current-project.name")
	}

	if projectName == "" {
		fmt.Println("请先使用【isx choose】选择项目")
		os.Exit(1)
	}

	if projectName == "isx-cli" {
		fmt.Println("该项目" + projectName + "暂不支持")
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

	// 执行 gradle docker 命令
	fmt.Printf("正在构建 %s 项目的Docker镜像...\n", projectName)

	gradleCmd := exec.Command("./gradlew", "docker")
	gradleCmd.Dir = projectPath
	gradleCmd.Stdout = os.Stdout
	gradleCmd.Stderr = os.Stderr

	err := gradleCmd.Run()
	if err != nil {
		fmt.Printf("Docker镜像构建失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Docker镜像构建成功！")
}
