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
	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: printCommand("isx build") + "| 编译项目代码",
	Long:  `isx build`,
	Run: func(cmd *cobra.Command, args []string) {
		buildCmdMain()
	},
}

func buildCmdMain() {

	projectName := viper.GetString("current-project.name")
	projectPath := viper.GetString(projectName+".dir") + "/" + viper.GetString(projectName+".name")
	buildImage := "registry.cn-shanghai.aliyuncs.com/isxcode/zhiqingyun-build"
	usr, _ := user.Current()

	// 获取gradle缓存目录
	cacheGradleDir := viper.GetString("cache.gradle.dir")
	if cacheGradleDir == "" {
		cacheGradleDir = usr.HomeDir + "/.gradle"
		viper.Set("cache.gradle.dir", cacheGradleDir)
		viper.WriteConfig()
	}
	_, err := os.Stat(cacheGradleDir)
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
		log.Fatal(err)
		fmt.Println("代码编译失败")
		os.Exit(1)
	} else {
		fmt.Println("代码编译完成")
	}
}
