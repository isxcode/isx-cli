package cmd

import (
	"fmt"
	"github.com/isxcode/isx-cli/common"
	"github.com/isxcode/isx-cli/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	Short: printCommand("isx clone", 65) + "| 下载项目代码",
	Long:  `isx clone`,
	Run: func(cmd *cobra.Command, args []string) {
		cloneCmdMain()
	},
}

func cloneCmdMain() {

	// 判断用户是否登录
	isLogin := common.CheckUserAccount(common.GetToken())
	if !isLogin {
		fmt.Println("请先登录")
		os.Exit(1)
	}

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
		fmt.Println("[" + strconv.Itoa(index) + "] " + printCommand(viper.GetString(projectName+".name"), 14) + "[ " + printCommand(viper.GetString(projectName+".repository.url"), 45) + "] : " + viper.GetString(projectName+".describe"))
	}

	// 输入项目编号
	fmt.Print("请输入下载项目编号：")
	fmt.Scanln(&projectNumber)
	projectName = projectList[projectNumber]
}

func inputProjectPath() {
	currentWorkDir := common.CurrentWorkDir()
	// 输入安装路径
	fmt.Printf("是否安装在当前路径(%s)下? (y/n) ", common.WhiteText(currentWorkDir))
	var flag = ""
	fmt.Scanln(&flag)
	flag = strings.Trim(flag, " ")
	flag = strings.ToUpper(flag)
	if flag != "Y" && flag != "N" {
		fmt.Println("输入值异常")
		os.Exit(1)
	}
	if flag == "Y" {
		projectPath = currentWorkDir
	}
	if flag == "N" {
		fmt.Print("请输入安装路径:")
		fmt.Scanln(&projectPath)
	}

	// 支持克隆路径替换～为当前用户目录
	if strings.HasPrefix(projectPath, "~/") {
		projectPath = strings.Replace(projectPath, "~", common.HomeDir(), 1)
	}
	projectPath = strings.ReplaceAll(projectPath, "\\", "/")

	// 目录不存在则报错
	_, err := os.Stat(projectPath)
	if os.IsNotExist(err) {
		fmt.Println("目录不存在，请重新输入")
		os.Exit(1)
	}

	// 目录不存在则报错
	_, err = os.Stat(projectPath + "/" + projectName)
	if err == nil {
		fmt.Println("项目已存在，请重新选择目录")
		os.Exit(1)
	}
}

func cloneCode(isxcodeRepository string, path string, name string, isMain bool) {

	// 替换下载链接
	isxcodeRepository = strings.Replace(isxcodeRepository, "https://", "https://"+common.GetToken()+"@", -1)

	// 下载主项目代码
	executeCommand := "git clone -b main " + isxcodeRepository
	cloneCmd := exec.Command("bash", "-c", executeCommand)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr
	cloneCmd.Dir = path
	err := cloneCmd.Run()
	if err != nil {
		fmt.Println("执行失败:", err)
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
	fmt.Println(updateOriginCommand)
	updateOriginCmd := exec.Command("bash", "-c", updateOriginCommand)
	updateOriginCmd.Stdout = os.Stdout
	updateOriginCmd.Stderr = os.Stderr
	updateOriginCmd.Dir = path + "/" + name
	updateOriginCmd.Run()

	// 添加upstream仓库
	addUpstreamCommand := "git remote add upstream " + isxcodeRepository + " && git fetch upstream"
	fmt.Println(addUpstreamCommand)
	addUpstreamCmd := exec.Command("bash", "-c", addUpstreamCommand)
	addUpstreamCmd.Stdout = os.Stdout
	addUpstreamCmd.Stderr = os.Stderr
	addUpstreamCmd.Dir = path + "/" + name
	addUpstreamCmd.Run()

	// main分支映射到isxcode仓库中
	linkUpstreamCommand := "git branch --set-upstream-to=upstream/main main"
	fmt.Println(linkUpstreamCommand)
	linkUpstreamCmd := exec.Command("bash", "-c", linkUpstreamCommand)
	linkUpstreamCmd.Stdout = os.Stdout
	linkUpstreamCmd.Stderr = os.Stderr
	linkUpstreamCmd.Dir = path + "/" + name
	linkUpstreamCmd.Run()
}

func cloneProjectCode() {

	// 下载主项目代码
	mainRepository := viper.GetString(projectName + ".repository.url")
	if !github.IsRepoForked(viper.GetString("user.account"), projectName) {
		github.ForkRepository("isxcode", projectName, "")
		cloneCode(mainRepository, projectPath, projectName, true)
	} else {
		cloneCode(mainRepository, projectPath, projectName, true)
	}

	// 下载子项目代码
	var subRepository []Repository
	viper.UnmarshalKey(projectName+".sub-repository", &subRepository)
	for _, repository := range subRepository {
		if !github.IsRepoForked(viper.GetString("user.account"), repository.Name) {
			forkRepository := github.ForkRepository("isxcode", repository.Name, "")
			if forkRepository {
				cloneCode(repository.Url, projectPath+"/"+repository.Name, repository.Name, false)
			}
		} else {
			cloneCode(repository.Url, projectPath+"/"+repository.Name, repository.Name, false)
		}
	}
}

func saveConfig() {
	viper.Set(projectName+".dir", projectPath)
	viper.Set("current-project.name", projectName)
	viper.WriteConfig()
}
