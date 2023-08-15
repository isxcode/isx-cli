package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(setCmd)
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: printCommand("isx set <config_key> <value>") + "| 设置配置参数",
	Long:  `isx set user.account ispong`,
	Run: func(cmd *cobra.Command, args []string) {
		setCmdMain(args)
	},
}

func setCmdMain(args []string) {
	viper.Set(args[0], args[1])
	viper.WriteConfig()
	fmt.Println("设置成功")
}
