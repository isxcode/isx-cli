package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: printCommand("isx version", 40) + "| 查看脚手架版本",
	Long:  `isx version`,
	Run: func(cmd *cobra.Command, args []string) {
		versionCmdMain()
	},
}

func versionCmdMain() {
	fmt.Println("当前版本号：v" + viper.GetString("version"))
}
