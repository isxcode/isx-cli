package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strconv"
)

func init() {
	rootCmd.AddCommand(chooseCmd)
}

var chooseCmd = &cobra.Command{
	Use:   "choose",
	Short: printCommand("isx choose") + "| 选择开发项目",
	Long:  `从isxcode组织中选择开发项目,isx choose`,
	Run: func(cmd *cobra.Command, args []string) {
		chooseCmdMain()
	},
}

func chooseCmdMain() {

	// 打印项目列表
	projectList := viper.GetStringSlice("project-list")
	for index, chooseProjectName := range projectList {
		fmt.Println("[" + strconv.Itoa(index) + "] " + viper.GetString(chooseProjectName+".name") + ": " + viper.GetString(chooseProjectName+".describe") + " 下载状态 【" + viper.GetString(chooseProjectName+".repository.download") + "】")
	}

	// 选择项目编号
	fmt.Print("请输入下载项目编号：")
	fmt.Scanln(&projectNumber)

	// 判断项目是否可切换
	projectName := projectList[projectNumber]
	isDownload := viper.GetString(projectName + ".repository.download")
	if isDownload != "ok" {
		fmt.Println("不可选择，请先下载代码")
		os.Exit(1)
	}

	// 设置当前的项目
	fmt.Println("切换到项目：" + projectName)
	viper.Set("current-project.name", projectName)
	viper.WriteConfig()
}
