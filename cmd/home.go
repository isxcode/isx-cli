package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(homeCmd)
}

var homeCmd = &cobra.Command{
	Use:   "home",
	Short: printCommand("isx home") + "| 快速进入项目目录",
	Long:  `isx home`,
	Run: func(cmd *cobra.Command, args []string) {
		homeCmdMain()
	},
}

func homeCmdMain() {
	projectName := viper.GetString("current-project.name")
	projectPath := viper.GetString(projectName+".dir") + "/" + projectName
	fmt.Println(projectPath)
}
