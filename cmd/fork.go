package cmd

import (
	"fmt"
	"github.com/isxcode/isx-cli/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var testExistFlag bool

func init() {
	forkCmd.Flags().BoolVarP(&testExistFlag, "test", "t", false, "测试是否已经fork过")
	rootCmd.AddCommand(forkCmd)
}

var forkCmd = &cobra.Command{
	Use:   "fork",
	Short: printCommand("isx fork", 65) + "| Fork当前项目为同名个人仓库",
	Long:  `isx fork | isx fork <project-name> | isx fork -t | isx fork -t <project-name>`,
	Run: func(cmd *cobra.Command, args []string) {
		ForkCmdMain(args)
	},
}

func ForkCmdMain(args []string) {
	if testExistFlag {
		if len(args) > 0 {
			forked := github.IsRepoForked(viper.GetString("user.account"), args[0])
			if forked {
				fmt.Println(args[0], "is forked!")
			} else {
				fmt.Println(args[0], "is not forked!")
			}
		} else {
			for _, project := range viper.GetStringSlice("project-list") {
				forked := github.IsRepoForked(viper.GetString("user.account"), viper.GetString(project+".name"))
				if forked {
					fmt.Println(viper.GetString(project+".name"), "is forked!")
				} else {
					fmt.Println(viper.GetString(project+".name"), "is not forked!")
				}
			}
		}
		return
	} else {
		projectName := viper.GetString("current-project.name")
		if len(args) > 0 {
			projectName = args[0]
		}
		github.ForkRepository("isxcode", projectName, "")
	}
}
