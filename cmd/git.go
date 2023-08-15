package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/exec"
)

var gitProjectPath string
var gitProjectName string

func init() {
	rootCmd.AddCommand(gitCmd)
}

var gitCmd = &cobra.Command{
	Use:   "git",
	Short: printCommand("isx git <git command>") + "| 子父项目执行git命令",
	Long:  `isx git <git command>`,
	Run: func(cmd *cobra.Command, args []string) {
		gitCmdMain(args)
	}, DisableFlagParsing: true,
}

func gitCmdMain(args []string) {

	gitProjectName = viper.GetString("current-project.name")
	gitProjectPath = viper.GetString(gitProjectName + ".dir")

	// 进入主项目执行git命令
	gitCmd := exec.Command("git", args...)
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr
	gitCmd.Dir = gitProjectPath + "/" + gitProjectName
	err := gitCmd.Run()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	} else {
		fmt.Println(gitProjectName + "git命令执行成功")
	}

	// 进入子项目执行命令
	var subRepository []Repository
	viper.UnmarshalKey(gitProjectName+".sub-repository", &subRepository)
	for _, repository := range subRepository {

		gitCmd := exec.Command("git", args...)
		gitCmd.Stdout = os.Stdout
		gitCmd.Stderr = os.Stderr
		gitCmd.Dir = gitProjectPath + "/" + gitProjectName + "/" + repository.Name
		err := gitCmd.Run()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		} else {
			fmt.Println(repository.Name + "git命令执行成功")
		}
	}
}
