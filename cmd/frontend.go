package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/exec"
	"runtime"
)

func init() {
	rootCmd.AddCommand(webCmd)
	rootCmd.AddCommand(frontendCmd)
}

var frontendCmd = &cobra.Command{
	Use:   "frontend",
	Short: printCommand("isx frontend", 40) + "| 本地启动前端",
	Long:  `isx frontend`,
	Run: func(cmd *cobra.Command, args []string) {
		webCmdMain()
	},
}

var webCmd = &cobra.Command{
	Use:    "frontend",
	Short:  printCommand("isx frontend", 40) + "| 本地启动前端服务",
	Long:   `isx frontend`,
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		webCmdMain()
	},
}

func webCmdMain() {
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
		projectDir := viper.GetString(projectName + ".dir")
		if projectDir != "" {
			projectPath = projectDir + "/" + projectName
		}
	}

	if projectPath == "" {
		fmt.Printf("项目 %s 未下载，请先使用【isx clone】下载项目代码\n", projectName)
		os.Exit(1)
	}

	var gradleCmd *exec.Cmd
	if runtime.GOOS == "windows" {
		gradleCmd = exec.Command("bash", "-c", "./gradlew.bat frontend")
	} else {
		gradleCmd = exec.Command("./gradlew", "frontend")
	}
	gradleCmd.Stdout = os.Stdout
	gradleCmd.Stderr = os.Stderr
	gradleCmd.Dir = projectPath
	err := gradleCmd.Run()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	} else {
		fmt.Println("执行成功")
	}
}
