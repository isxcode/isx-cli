package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

var deleteProjectNumber int

func init() {
	rootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: printCommand("isx remove", 65) + "| 删除本地项目",
	Long:  `isx remove`,
	Run: func(cmd *cobra.Command, args []string) {
		removeCmdMain()
	},
}

func removeCmdMain() {

	// 选择项目编号
	inputRemoveProjectNumber()

	// 删除项目
	removeProject()
}

func inputRemoveProjectNumber() {

	// 打印项目列表
	projectList := viper.GetStringSlice("project-list")
	for index, projectName := range projectList {
		fmt.Println("[" + strconv.Itoa(index) + "] " + printCommand(viper.GetString(projectName+".name"), 14) + "[ " + printCommand(viper.GetString(projectName+".repository.url"), 45) + "] : " + viper.GetString(projectName+".describe"))
	}

	// 输入项目编号
	fmt.Print("请输入删除项目编号：")
	fmt.Scanln(&deleteProjectNumber)

	// 没有下载的不让删除
	projectName := projectList[deleteProjectNumber]
	downloadStatus := viper.GetString(projectName + ".repository.download")
	if downloadStatus != "ok" {
		fmt.Print("该项目未下载")
		os.Exit(1)
	}

	// 二次确认
	deleteProject := "N"
	fmt.Print("确认要删除该项目吗？(Y/N) default is N: ")
	fmt.Scanln(&deleteProject)
	if deleteProject == "N" {
		fmt.Println("已中止")
		os.Exit(1)
	}
}

func removeProject() {

	// 获取项目目录
	projectList := viper.GetStringSlice("project-list")
	projectName := projectList[deleteProjectNumber]
	projectPath := viper.GetString(projectName + ".dir")

	// 三次确认删除
	deleteProject := "N"
	fmt.Print("确认要删除【" + projectPath + "/" + projectName + "】路径吗?(Y/N) default is N: ")
	fmt.Scanln(&deleteProject)
	if deleteProject == "N" {
		fmt.Println("已中止")
		os.Exit(1)
	}

	// 更新平台替换projectPath
	removeCommand := ""
	if runtime.GOOS == "windows" {
		projectPath = strings.ReplaceAll(projectPath, "C:", "/c")
		projectPath = strings.ReplaceAll(projectPath, " ", "\\ ")
		removeCommand = "rm -rf " + projectPath + "/" + projectName
	} else {
		removeCommand = "rm -rf " + projectPath + "/" + projectName
	}

	removeCmd := exec.Command("bash", "-c", removeCommand)
	removeCmd.Stdout = os.Stdout
	removeCmd.Stderr = os.Stderr
	err := removeCmd.Run()
	if err != nil {
		fmt.Println("执行失败:", err)
		os.Exit(1)
	} else {
		fmt.Println(projectPath + "/" + projectName + "路径已删除")
	}

	// 保存配置
	if viper.GetString("current-project.name") == projectName {
		viper.Set("current-project.name", "")
	}
	viper.Set(projectName+".dir", "")
	viper.Set(projectName+".repository.download", "no")
	viper.WriteConfig()
}
