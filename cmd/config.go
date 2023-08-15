/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: printCommand("isx config") + "| 查看配置文件",
	Long:  `查看配置内容，举例：isx config`,
	Run: func(cmd *cobra.Command, args []string) {
		data, err := yaml.Marshal(viper.AllSettings())
		if err != nil {
			fmt.Println("Failed to marshal projects:", err)
			return
		}
		fmt.Println(string(data))
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
