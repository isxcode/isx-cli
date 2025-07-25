package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/exec"
	"os/user"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:    "run",
	Short:  printCommand("isx run [frontend/backend/web] [port]", 40) + "| 使用docker运行项目",
	Long:   `isx run frontend 8888/ isx run backend 8888/ isx run 8888/isx run web 8888`,
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {

		runType := ""
		port := ""

		if len(args) < 1 {
			runType = "all"
			port = "8080"
		}

		if len(args) == 1 {
			switch args[0] {
			case "backend":
				runType = args[0]
				port = "8080"
				break
			case "frontend":
				runType = args[0]
				port = "5173"
				break
			case "web":
				runType = args[0]
				port = "3000"
				break
			default:
				runType = "all"
				port = args[0]
			}
		}

		if len(args) == 2 {
			runType = args[0]
			port = args[1]
		}

		if len(args) > 2 {
			fmt.Println("使用方式不对，请重新输入命令")
			os.Exit(1)
		}

		if runType != "all" && runType != "backend" && runType != "web" && runType != "frontend" {
			fmt.Println("使用方式不对，请重新输入命令")
			os.Exit(1)
		}

		runCmdMain(runType, port)
	},
}

func runCmdMain(runType string, port string) {
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
		fmt.Println("当前项目未下载，请先使用 【isx clone】 下载项目代码")
		os.Exit(1)
	}

	usr, _ := user.Current()

	// 获取gradle缓存目录
	cacheGradleDir := viper.GetString("cache.gradle.dir")
	if cacheGradleDir == "" {
		cacheGradleDir = usr.HomeDir + "/.gradle"
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
		cachePnpmDir = usr.HomeDir + "/.pnpm-store"
	}
	_, err = os.Stat(cachePnpmDir)
	if os.IsNotExist(err) {
		err := os.Mkdir(cachePnpmDir, 0755)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	var runCommand string
	switch runType {
	case "backend":
		runCommand = "docker run " +
			"-v " + projectPath + ":/spark-yun " +
			"-v " + cacheGradleDir + ":/root/.gradle " +
			"-e ENV_TYPE='BACKEND' " +
			"-e TZ=Asia/Shanghai " +
			"-p " + port + ":8080 " +
			" registry.cn-shanghai.aliyuncs.com/isxcode/zhiqingyun-local:latest"
		break
	case "frontend":
		runCommand = "docker run " +
			"-v " + projectPath + ":/spark-yun " +
			"-v " + cachePnpmDir + ":/root/.pnpm-store " +
			"-e ENV_TYPE='FRONTEND' " +
			"-e TZ=Asia/Shanghai " +
			"-p " + port + ":5173 " +
			" registry.cn-shanghai.aliyuncs.com/isxcode/zhiqingyun-local:latest"
		break
	case "all":
		runCommand = "docker run " +
			"-v " + projectPath + ":/spark-yun " +
			"-e ENV_TYPE='ALL' " +
			"-e TZ=Asia/Shanghai " +
			"-p " + port + ":8080 " +
			" registry.cn-shanghai.aliyuncs.com/isxcode/zhiqingyun-local:latest"
		break
	case "web":
		runCommand = "docker run " +
			"-e ENV_TYPE='WEB' " +
			"-e TZ=Asia/Shanghai " +
			"-p " + port + ":3000 " +
			"-v " + projectPath + ":/spark-yun " +
			"-v " + cachePnpmDir + ":/root/.pnpm-store " +
			" registry.cn-shanghai.aliyuncs.com/isxcode/zhiqingyun-local:latest"
		break
	default:
		fmt.Println("使用方式不对，请重新输入命令")
		os.Exit(1)
	}

	runCmd := exec.Command("bash", "-c", runCommand)
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr
	err = runCmd.Run()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	} else {
		fmt.Println("代码运行完毕")
	}
}
