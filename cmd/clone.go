package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
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
	isLogin := checkUserAccount()
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

func checkUserAccount() bool {

	headers := http.Header{}
	headers.Set("Accept", "application/vnd.github+json")
	headers.Set("Authorization", "Bearer "+viper.GetString("user.token"))
	headers.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/octocat", nil)

	req.Header = headers
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败:", err)
		os.Exit(1)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// 解析结果
	if resp.StatusCode == http.StatusOK {
		return true
	} else {
		return false
	}
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

	// 输入安装路径
	fmt.Print("请输入安装路径:")
	fmt.Scanln(&projectPath)

	// 支持克隆路径替换～为当前用户目录
	if strings.HasPrefix(projectPath, "~/") {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		projectPath = strings.Replace(projectPath, "~", home, 1)
	}

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
