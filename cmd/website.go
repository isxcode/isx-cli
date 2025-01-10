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
	rootCmd.AddCommand(websiteCmd)
}

var websiteCmd = &cobra.Command{
	Use:   "website",
	Short: printCommand("isx website", 65) + "| 本地启动官网",
	Long:  `isx website`,
	Run: func(cmd *cobra.Command, args []string) {
		websiteCmdMain()
	},
}

func websiteCmdMain() {

	projectName := viper.GetString("current-project.name")
	projectDir := viper.GetString(projectName + ".dir")
	projectPath := projectDir + "/" + projectName

	var gradleCmd *exec.Cmd
	if runtime.GOOS == "windows" {
		gradleCmd = exec.Command("bash", "-c", "./gradlew.bat website")
	} else {
		gradleCmd = exec.Command("./gradlew", "website")
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
