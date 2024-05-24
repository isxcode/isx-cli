/*
Copyright © 2024 jamie HERE <EMAIL ADDRESS>
*/
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
	rootCmd.AddCommand(dbCmd)
}

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: printCommand("isx db list | isx db <issue_number>", 65) + "| 查看当前db(暂不开放)",
	Long:  `isx db | isx db list | isx db 123`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Println("使用方式不对，请重新输入命令")
		}
		dbCmdMain(args)
	},
}

func dbCmdMain(args []string) {
	projectName := viper.GetString("current-project.name")
	if projectName != "spark-yun" {
		fmt.Println("目前仅 spark-yun 项目支持该命令")
		os.Exit(1)
	}

	if len(args) == 0 || args[0] == "list" {
		listDb()
		os.Exit(0)
	}

	// TODO 此处应为项目名称 例 spark-yun
	projectName = ".zhiqingyun"
	createDb(projectName, args[0])
	checkoutDb(projectName, args[0])
}

func createDb(projectName, issueNumber string) {
	dbDir := generatePath(projectName, issueNumber)
	os.MkdirAll(dbDir, 0755)
}

func checkoutDb(projectName, issueNumber string) {
	// TODO diDir为 项目启动时 使用的h2路径 若多项目公用同一引用路径 可修改generatePath方法
	// 此处目前为 /Users/jamie/.zhiqingyun/h2
	dbDir := generatePath(projectName, "")
	// 删除文件 当删除文件夹时会报错
	rmCmd := exec.Command("rm", "-f", dbDir)
	rmCmd.Stdout = os.Stdout
	rmCmd.Stderr = os.Stderr
	err := rmCmd.Run()
	if err != nil {
		os.Exit(1)
	}
	// 软链接
	lnCmd := exec.Command("ln", "-s", generatePath(projectName, issueNumber), dbDir)
	err = lnCmd.Run()
	if err != nil {
		fmt.Println("Failed to run ln command:", err)
	}
	fmt.Println("db 切换成功")
}

func listDb() {
	// TODO 目录待优化
	lsCmd := exec.Command("ls", common.HomeDir()+"/.zhiqingyun")
	output, _ := lsCmd.Output()
	dirs := string(output)
	split := strings.Split(dirs, "\n")
	for i, dir := range split[:len(split)-1] {
		if dir == "h2" {
			continue
		}
		// 输出格式待优化
		fmt.Printf("[%2d]: %s \n", i, dir)
	}

}

func generatePath(projectName, issueNumber string) string {
	// TODO 统一路径 例 common.HomeDir() + “/” + “.isx”
	projectHome := common.HomeDir()

	basePath := projectHome + "/" + projectName + "/"
	if issueNumber == "" {
		return basePath + "h2"
	}
	return basePath + issueNumber + "/" + "h2"
}
