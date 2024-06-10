package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"runtime"
)

func init() {
	rootCmd.AddCommand(formatCmd)
}

var formatCmd = &cobra.Command{
	Use:   "format",
	Short: printCommand("isx format", 65) + "| 格式化代码",
	Long:  `isx format`,
	Run: func(cmd *cobra.Command, args []string) {
		formatCmdMain()
	},
}

func formatCmdMain() {

	projectName := viper.GetString("current-project.name")
	projectPath := viper.GetString(projectName+".dir") + "/" + projectName

	// 除了isx-cli项目，其他都要使用gradle 格式化代码
	if "isx-cli" != projectName && "tools-yun" != projectName {

		var gradleCmd *exec.Cmd
		if runtime.GOOS == "windows" {
			gradleCmd = exec.Command("bash", "-c", "./gradlew.bat format")
		} else {
			gradleCmd = exec.Command("./gradlew", "format")
		}
		gradleCmd.Stdout = os.Stdout
		gradleCmd.Stderr = os.Stderr
		gradleCmd.Dir = projectPath
		err := gradleCmd.Run()
		if err != nil {
			fmt.Println("执行失败:", err)
			os.Exit(1)
		} else {
			fmt.Println("执行成功")
		}
	}
}
