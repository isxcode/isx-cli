package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func init() {
	rootCmd.AddCommand(backendCmd)
}

var backendCmd = &cobra.Command{
	Use:   "backend",
	Short: printCommand("isx backend", 65) + "| 使用gradle启动项目",
	Long:  `isx backend`,
	Run: func(cmd *cobra.Command, args []string) {
		backendCmdMain()
	},
}

func backendCmdMain() {

	projectName := viper.GetString("current-project.name")
	projectDir := viper.GetString(projectName + ".dir")
	projectPath := filepath.Join(projectDir, projectName)

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
