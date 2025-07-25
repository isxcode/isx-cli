package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/exec"
)

func init() {
	rootCmd.AddCommand(uploadCmd)
}

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: printCommand("isx upload", 40) + "| 发布本地安装包",
	Long:  `上传项目的构建产物到仓库`,
	Run: func(cmd *cobra.Command, args []string) {
		uploadCmdMain()
	},
}

func uploadCmdMain() {
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

	// 执行 gradle upload 命令
	fmt.Printf("正在上传 %s 项目的构建产物...\n", projectName)

	gradleCmd := exec.Command("./gradlew", "upload")
	gradleCmd.Dir = projectPath
	gradleCmd.Stdout = os.Stdout
	gradleCmd.Stderr = os.Stderr

	err := gradleCmd.Run()
	if err != nil {
		fmt.Printf("构建产物上传失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("构建产物上传成功！")
}
