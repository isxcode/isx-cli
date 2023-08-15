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
	rootCmd.AddCommand(formatCmd)
}

var formatCmd = &cobra.Command{
	Use:   "format",
	Short: printCommand("isx format") + "| 代码格式化",
	Long:  `isx format`,
	Run: func(cmd *cobra.Command, args []string) {
		formatCmdMain()
	},
}

func formatCmdMain() {

	projectName := viper.GetString("current-project.name")
	projectPath := viper.GetString(projectName+".dir") + "/" + projectName

	gradleCmd := exec.Command("./gradlew", "format")
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
