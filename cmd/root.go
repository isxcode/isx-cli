/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"github.com/isxcode/isx-cli/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "isx",
	Short: "cli for isxcode app",
	Long: `
 ____ _____ __ __           __  _      ____ 
|    / ___/|  |  |         /  ]| |    |    |
 |  (   \_ |  |  | _____  /  / | |     |  | 
 |  |\__  ||_   _||     |/  /  | |___  |  | 
 |  |/  \ ||     ||_____/   \_ |     | |  | 
 |  |\    ||  |  |      \     ||     | |  | 
|____|\___||__|__|       \____||_____||____|

至行云-至爻数据开发规范脚手架
代码仓库：https://github.com/isxcode/isx-cli

`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cmd.DisableFlagParsing = true
	},

	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var (
	cfgFile string
)

type Repository struct {
	Download string `yaml:"download"`
	Url      string `yaml:"url"`
	Name     string `yaml:"name"`
}

func init() {

	cobra.OnInitialize(initConfig)

	// 禁用自动生成的 completion 命令
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// 禁用自动生成的 help 命令
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})

	// 隐藏使用说明和标志
	rootCmd.SetUsageTemplate(`{{if .Runnable}}{{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}

`)

	// 解析配置文件
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.isx/config.yml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {

	// 获取home目录
	home := common.HomeDir()

	// 初始化配置文件信息
	viper.SetConfigFile(home + "/.isx/config.yml")

	// 判断配置文件是否存在
	if err := viper.ReadInConfig(); err != nil {

		// 判断文件夹是否存在，不存在则新建
		_, err := os.Stat(home + "/.isx")
		if os.IsNotExist(err) {
			err := os.Mkdir(home+"/.isx", 0755)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		// 判断文件是否存在，不存在则新建
		_, err = os.Stat(home + "/.isx/config.yml")
		if os.IsNotExist(err) {
			// 初始化配置
			viper.SetConfigType("yaml")

			// 生成随机密钥用于加密
			encryptionKey := common.GenerateEncryptionKey()

			var yamlExample = []byte(`version: 1.1.2
project-list:
    - name: spark-yun
      describe: 至轻云-超轻量级智能化数据中心
      repository-url: https://github.com/isxcode/spark-yun.git
      dir: ""
      sub-repository:
        - name: spark-yun-vip
          url: https://github.com/isxcode/spark-yun-vip.git
    - name: torch-yun
      describe: 至数云-超轻量级一体化应用平台
      repository-url: https://github.com/isxcode/torch-yun.git
      dir: ""
      sub-repository:
        - name: torch-yun-vip
          url: https://github.com/isxcode/torch-yun-vip.git
    - name: isx-cli
      describe: 至行云-至爻数据开发规范脚手架
      repository-url: https://github.com/isxcode/isx-cli.git
      dir: ""
now-project: ""
user:
    account: ""
    token: ""
    secret: ` + encryptionKey + `
`)
			err := viper.ReadConfig(bytes.NewBuffer(yamlExample))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			// 持久化配置
			err = viper.SafeWriteConfigAs(home + "/.isx/config.yml")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}
}

func printCommand(commandDesc string, length int) string {

	return commandDesc + strings.Repeat(" ", length-len(commandDesc))
}
