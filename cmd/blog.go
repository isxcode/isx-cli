package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const blogProjectName = "ispong-blogs"
const blogRepositoryOwner = "ispong"
const blogRepositoryURL = "https://github.com/ispong/ispong-blogs.git"

func init() {
	blogCmd.AddCommand(blogPushCmd)
	blogCmd.AddCommand(blogPullCmd)
	rootCmd.AddCommand(blogCmd)
}

var blogCmd = &cobra.Command{
	Use:   "blog",
	Short: printCommand("isx blog <command>", 40) + "| 博客管理",
	Long:  `isx blog <command>`,
}

var blogPushCmd = &cobra.Command{
	Use:   "push",
	Short: printCommand("isx blog push", 40) + "| 提交博客",
	Long:  `isx blog push`,
	Run: func(cmd *cobra.Command, args []string) {
		blogPushCmdMain()
	},
}

var blogPullCmd = &cobra.Command{
	Use:   "pull",
	Short: printCommand("isx blog pull", 40) + "| 拉取博客",
	Long:  `isx blog pull`,
	Run: func(cmd *cobra.Command, args []string) {
		blogPullCmdMain()
	},
}

func getBlogDir() string {
	blogDir := viper.GetString("blog.dir")
	if blogDir == "" {
		fmt.Println("请先使用【isx clone】下载博客项目，或使用【isx set blog.dir <blog_path>】设置博客目录")
		os.Exit(1)
	}
	return blogDir
}

func checkBlogProject() string {
	blogDir := getBlogDir()

	if _, err := os.Stat(blogDir); os.IsNotExist(err) {
		fmt.Println("博客目录不存在：" + blogDir)
		fmt.Println("请先使用【isx clone】下载博客项目，或使用【isx set blog.dir <blog_path>】设置博客目录")
		os.Exit(1)
	}

	if _, err := os.Stat(blogDir + "/package.json"); os.IsNotExist(err) {
		fmt.Println("当前目录不是有效的Hexo博客项目：" + blogDir)
		fmt.Println("请确认【blog.dir】配置的是博客项目根目录")
		os.Exit(1)
	}

	return blogDir
}

func checkBlogHexoDependencies() string {
	blogDir := checkBlogProject()

	if _, err := os.Stat(blogDir + "/node_modules/.bin/hexo"); os.IsNotExist(err) {
		fmt.Println("博客项目Hexo依赖未安装")
		fmt.Println("开始自动安装博客依赖")
		installBlogDependencies(blogDir)
	}

	return blogDir
}

func blogPushCmdMain() {
	blogDir := checkBlogHexoDependencies()

	commands := [][]string{
		{"git", "add", "."},
		{"git", "commit", "-m", ":memo: 写博客"},
		{"git", "push", "origin", "main"},
	}

	for _, command := range commands {
		cmd := exec.Command(command[0], command[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = blogDir
		err := cmd.Run()
		if err != nil {
			fmt.Println("执行失败:", err)
			os.Exit(1)
		}
	}

	fmt.Println("提交博客成功")
}

func blogPullCmdMain() {
	blogDir := checkBlogProject()

	commands := [][]string{
		{"git", "fetch", "origin"},
		{"git", "rebase", "origin/main"},
	}

	for _, command := range commands {
		cmd := exec.Command(command[0], command[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = blogDir
		err := cmd.Run()
		if err != nil {
			fmt.Println("执行失败:", err)
			os.Exit(1)
		}
	}

	fmt.Println("拉取博客成功")
}
