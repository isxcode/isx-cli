package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/exec"
)

func init() {
	rootCmd.AddCommand(gradleCmd)
}

var gradleCmd = &cobra.Command{
	Use:   "gradle",
	Short: printCommand("isx gradle") + "| 执行项目gradle命令",
	Long:  `gradle install、gradle start、gradle clean、gradle format`,
	Run: func(cmd *cobra.Command, args []string) {
		gradleCmdMain(args)
	},
}

func gradleCmdMain(args []string) {

	projectName := viper.GetString("current-project.name")
	projectPath := viper.GetString(projectName+".dir") + "/" + projectName

	gradleCmd := exec.Command("./gradlew", args...)
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
