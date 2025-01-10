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
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: printCommand("isx server", 65) + "| 本地启动后端程序",
	Long:  `isx server`,
	Run: func(cmd *cobra.Command, args []string) {
		serverCmdMain()
	},
}

func serverCmdMain() {

	projectName := viper.GetString("current-project.name")
	projectDir := viper.GetString(projectName + ".dir")
	projectPath := projectDir + "/" + projectName

	var gradleCmd *exec.Cmd
	if runtime.GOOS == "windows" {
		gradleCmd = exec.Command("bash", "-c", "./gradlew.bat backend")
	} else {
		gradleCmd = exec.Command("./gradlew", "backend")
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
