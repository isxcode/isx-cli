package cmd

import (
	"fmt"
	"github.com/isxcode/isx-cli/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"strings"
)

func init() {
	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:    "build",
	Short:  printCommand("isx build", 40) + "| 使用docker编译项目代码",
	Long:   `isx build,大约需要10分钟,需要docker环境`,
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		buildCmdMain()
	},
}

func buildCmdMain() {
	// 定义项目结构体
	type ProjectConfig struct {
		Name          string `mapstructure:"name"`
		Describe      string `mapstructure:"describe"`
		RepositoryURL string `mapstructure:"repository-url"`
		Dir           string `mapstructure:"dir"`
	}

	projectName := viper.GetString("now-project")
	if projectName == "" {
		fmt.Println("没有选择当前项目，请先使用 'isx choose' 选择项目")
		os.Exit(1)
	}

	// 获取项目列表
	var projectList []ProjectConfig
	err := viper.UnmarshalKey("project-list", &projectList)
	if err != nil {
		fmt.Printf("读取项目列表失败: %v\n", err)
		os.Exit(1)
	}

	// 找到当前项目的路径
	var projectPath string
	for _, proj := range projectList {
		if proj.Name == projectName {
			projectPath = proj.Dir
			break
		}
	}

	if projectPath == "" {
		fmt.Println("当前项目未下载，请先使用 'isx clone' 下载项目代码")
		os.Exit(1)
	}

	buildImage := "registry.cn-shanghai.aliyuncs.com/isxcode/zhiqingyun-build"
	home := common.HomeDir()

	// 获取gradle缓存目录
	cacheGradleDir := viper.GetString("cache.gradle.dir")
	if cacheGradleDir == "" {
		cacheGradleDir = home + "/.gradle"
		cacheGradleDir = strings.ReplaceAll(cacheGradleDir, "\\", "/")
		viper.Set("cache.gradle.dir", cacheGradleDir)
		viper.WriteConfig()
	}
	_, err = os.Stat(cacheGradleDir)
	if os.IsNotExist(err) {
		err := os.Mkdir(cacheGradleDir, 0755)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// 获取pnpm缓存目录
	cachePnpmDir := viper.GetString("cache.pnpm.dir")
	if cachePnpmDir == "" {
		cachePnpmDir = home + "/.pnpm-store"
		cachePnpmDir = strings.ReplaceAll(cachePnpmDir, "\\", "/")
		viper.Set("cache.pnpm.dir", cachePnpmDir)
		viper.WriteConfig()
	}
	_, err = os.Stat(cachePnpmDir)
	if os.IsNotExist(err) {
		err := os.Mkdir(cachePnpmDir, 0755)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// 镜像编译代码
	buildCommand := "docker run " +
		"--rm " +
		"-v " + projectPath + ":/spark-yun " +
		"-v " + cachePnpmDir + ":/root/.pnpm-store " +
		"-v " + cacheGradleDir + ":/root/.gradle " +
		buildImage
	buildCmd := exec.Command("bash", "-c", buildCommand)
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	err = buildCmd.Run()
	if err != nil {
		fmt.Println("代码编译失败", err)
		os.Exit(1)
	} else {
		fmt.Println("代码编译完成")
	}
}
