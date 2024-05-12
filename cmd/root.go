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

欢迎使用isx-cli脚手架
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

	// 解析配置文件
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.isx/isx-config.yml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {

	// 获取home目录
	home := common.HomeDir()

	// 初始化配置文件信息
	viper.SetConfigFile(home + "/.isx/isx-config.yml")

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
		_, err = os.Stat(home + "/.isx/isx-config.yml")
		if os.IsNotExist(err) {
			// 初始化配置
			viper.SetConfigType("yaml")
			var yamlExample = []byte(`
current-project:
    name: ""
cache:
    gradle:
        dir: ""
    pnpm:
        dir: ""
user:
    account: ""
    token: ""
project-list:
    - spark-yun
    - flink-yun
    - echarts-yun
    - pytorch-yun
    - meta-yun
    - isx-cli
spark-yun:
    name: spark-yun
    describe: 至轻云-打造大数据计算平台
    dir: ""
    repository:
        url: https://github.com/isxcode/spark-yun.git
        download: no
    sub-repository:
        - url: https://github.com/isxcode/spark-yun-vip.git
          name: spark-yun-vip
flink-yun:
    name: flink-yun
    describe: 至流云-打造流数据分析平台
    dir: ""
    repository:
        url: https://github.com/isxcode/flink-yun.git
        download: no
    sub-repository:
        - url: https://github.com/isxcode/flink-yun-vip.git
          name: flink-yun-vip
pytorch-yun:
    name: pytorch-yun
    describe: 至慧云-打造智能微模型平台
    dir: ""
    repository:
        url: https://github.com/isxcode/pytorch-yun.git
        download: no
    sub-repository:
        - url: https://github.com/isxcode/pytorch-yun-vip.git
          name: pytorch-yun-vip
meta-yun:
    name: meta-yun
    describe: 至数云-打造元数据管理平台
    dir: ""
    repository:
        url: https://github.com/isxcode/meta-yun.git
        download: no
    sub-repository:
        - url: https://github.com/isxcode/meta-yun-vip.git
          name: meta-yun-vip
echarts-yun:
    name: echarts-yun
    describe: 至视云-打造数据可视化平台
    dir: ""
    repository:
        url: https://github.com/isxcode/echarts-yun.git
        download: no
    sub-repository:
        - url: https://github.com/isxcode/echarts-yun-vip.git
          name: echarts-yun-vip
isx-cli:
    name: isx-cli
    describe: 至行云-打造开发规范脚手架
    dir: ""
    repository:
        url: https://github.com/isxcode/isx-cli.git
        download: no
    sub-repository:
version:
    number: 1.1.0
`)
			err := viper.ReadConfig(bytes.NewBuffer(yamlExample))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			// 持久化配置
			err = viper.SafeWriteConfigAs(home + "/.isx/isx-config.yml")
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
