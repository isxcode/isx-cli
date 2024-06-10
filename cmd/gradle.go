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
	"strings"
)

func init() {
	rootCmd.AddCommand(gradleCmd)
}

var gradleCmd = &cobra.Command{
	Use:   "gradle",
	Short: printCommand("isx gradle <gradle_command>", 65) + "| 项目内执行gradle命令",
	Long:  `isx gradle install、isx gradle start、isx gradle clean、isx gradle format`,
	Run: func(cmd *cobra.Command, args []string) {
		gradleCmdMain(args)
	},
}

func gradleCmdMain(args []string) {

	projectName := viper.GetString("current-project.name")
	projectDir := viper.GetString(projectName + ".dir")
	projectPath := filepath.Join(projectDir, projectName)

	var gradleCmd *exec.Cmd
	if runtime.GOOS == "windows" {
		gradleCmd = exec.Command("bash", "-c", "./gradlew.bat "+strings.Join(args, " "))
	} else {
		gradleCmd = exec.Command("./gradlew", args...)
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
