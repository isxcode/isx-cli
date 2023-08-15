package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var projectNumber int
var projectPath string
var projectName string

func init() {
	rootCmd.AddCommand(cloneCmd)
}

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: printCommand("isx clone") + "| 下载项目代码",
	Long:  `isx clone`,
	Run: func(cmd *cobra.Command, args []string) {
		cloneCmdMain()
	},
}

func cloneCmdMain() {

	// 选择项目编号
	inputProjectNumber()

	// 输入安装路径
	inputProjectPath()

	// 下载项目代码
	cloneProjectCode()

	// 保存配置
	saveConfig()
}

func inputProjectNumber() {

	// 打印项目列表
	projectList := viper.GetStringSlice("project-list")
	for index, projectName := range projectList {
		fmt.Println("[" + strconv.Itoa(index) + "] " + viper.GetString(projectName+".name") + ": " + viper.GetString(projectName+".describe"))
	}

	// 输入项目编号
	fmt.Print("请输入下载项目编号：")
	fmt.Scanln(&projectNumber)
	projectName = projectList[projectNumber]
}

func inputProjectPath() {

	// 输入安装路径
	fmt.Print("请输入安装路径:")
	fmt.Scanln(&projectPath)

	// 目录不存在则报错
	_, err := os.Stat(projectPath)
	if os.IsNotExist(err) {
		fmt.Println("目录不存在，请重新输入")
		os.Exit(1)
	}
}

func cloneCode(isxcodeRepository string, path string, name string, isMain bool) {

	// 替换下载链接
	isxcodeRepository = strings.Replace(isxcodeRepository, "https://", "https://"+viper.GetString("user.token")+"@", -1)

	// 下载主项目代码
	executeCommand := "git clone -b main " + isxcodeRepository
	cloneCmd := exec.Command("bash", "-c", executeCommand)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr
	cloneCmd.Dir = path
	err := cloneCmd.Run()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	} else {
		if isMain {
			viper.Set(projectName+".repository.download", "ok")
			viper.WriteConfig()
		}
		fmt.Println(name + "下载成功")
	}

	// 将origin改为个人的
	userRepository := strings.Replace(isxcodeRepository, "isxcode", viper.GetString("user.account"), -1)
	updateOriginCommand := "git remote set-url origin " + userRepository + " && git fetch origin"
	updateOriginCmd := exec.Command("bash", "-c", updateOriginCommand)
	updateOriginCmd.Stdout = os.Stdout
	updateOriginCmd.Stderr = os.Stderr
	updateOriginCmd.Dir = path + "/" + name
	updateOriginCmd.Run()

	// 添加upstream仓库
	addUpstreamCommand := "git remote add upstream " + isxcodeRepository + " && git fetch upstream"
	addUpstreamCmd := exec.Command("bash", "-c", addUpstreamCommand)
	addUpstreamCmd.Stdout = os.Stdout
	addUpstreamCmd.Stderr = os.Stderr
	addUpstreamCmd.Dir = path + "/" + name
	addUpstreamCmd.Run()
}

func cloneProjectCode() {

	// 下载主项目代码
	mainRepository := viper.GetString(projectName + ".repository.url")
	cloneCode(mainRepository, projectPath, projectName, true)

	// 下载子项目代码
	var subRepository []Repository
	viper.UnmarshalKey(projectName+".sub-repository", &subRepository)
	for _, repository := range subRepository {
		cloneCode(repository.Url, projectPath+"/"+projectName, repository.Name, false)
	}
}

func saveConfig() {
	viper.Set(projectName+".dir", projectPath)
	viper.Set("current-project.name", projectName)
	viper.WriteConfig()
}
