package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/exec"
	"runtime"
)

func init() {
	rootCmd.AddCommand(webCmd)
}

var webCmd = &cobra.Command{
	Use:    "web",
	Short:  printCommand("isx web", 65) + "| 本地启动前端服务",
	Long:   `isx web`,
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		webCmdMain()
	},
}

func webCmdMain() {

	projectName := viper.GetString("current-project.name")
	projectDir := viper.GetString(projectName + ".dir")
	projectPath := projectDir + "/" + projectName

	var gradleCmd *exec.Cmd
	if runtime.GOOS == "windows" {
		gradleCmd = exec.Command("bash", "-c", "./gradlew.bat frontend")
	} else {
		gradleCmd = exec.Command("./gradlew", "frontend")
	}
	gradleCmd.Stdout = os.Stdout
	gradleCmd.Stderr = os.Stderr
	gradleCmd.Dir = projectPath
	err := gradleCmd.Run()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	} else {
		fmt.Println("执行成功")
	}
}
