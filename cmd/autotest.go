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
	rootCmd.AddCommand(autotestCmd)
}

var autotestCmd = &cobra.Command{
	Use:   "autotest",
	Short: printCommand("isx autotest", 65) + "| 自动化测试",
	Long:  `isx autotest`,
	Run: func(cmd *cobra.Command, args []string) {
		autotestCmdMain()
	},
}

func autotestCmdMain() {

	projectName := viper.GetString("current-project.name")
	projectDir := viper.GetString(projectName + ".dir")
	projectPath := projectDir + "/" + projectName

	var gradleCmd *exec.Cmd
	if runtime.GOOS == "windows" {
		gradleCmd = exec.Command("bash", "-c", "./gradlew.bat autotest")
	} else {
		gradleCmd = exec.Command("./gradlew", "autotest")
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
