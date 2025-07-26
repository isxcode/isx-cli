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
	Short: printCommand("isx upload <target>", 40) + "| 发布本地安装包",
	Long:  `上传项目到指定仓库，支持：oss、docker、ali`,
	Run: func(cmd *cobra.Command, args []string) {
		uploadCmdMain(args)
	},
}

func uploadCmdMain(args []string) {
	// 检查参数数量
	if len(args) != 1 {
		fmt.Println("使用方式不对，请输入：isx upload <target>")
		fmt.Println("支持的目标：oss、docker、ali")
		os.Exit(1)
	}

	// 验证参数值
	target := args[0]
	var gradleTask string
	switch target {
	case "oss":
		gradleTask = "upload-ali-oss"
	case "docker":
		gradleTask = "upload-docker-hub"
	case "ali":
		gradleTask = "upload-ali-hub"
	default:
		fmt.Printf("不支持的目标：%s\n", target)
		fmt.Println("支持的目标：oss、docker、ali")
		os.Exit(1)
	}
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

	// 执行对应的 gradle 命令
	fmt.Printf("正在发布 %s 项目到 %s...\n", projectName, target)

	gradleCmd := exec.Command("./gradlew", gradleTask)
	gradleCmd.Dir = projectPath
	gradleCmd.Stdout = os.Stdout
	gradleCmd.Stderr = os.Stderr

	err := gradleCmd.Run()
	if err != nil {
		fmt.Printf("上传到 %s 失败: %v\n", target, err)
		os.Exit(1)
	}

	fmt.Printf("上传到 %s 成功！\n", target)
}
