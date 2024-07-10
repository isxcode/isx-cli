/*
Copyright © 2024 jamie HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/isxcode/isx-cli/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
)

func init() {
	rootCmd.AddCommand(forkCmd)
}

var forkCmd = &cobra.Command{
	Use:   "fork",
	Short: printCommand("isx fork", 65) + "| Fork当前项目为同名个人仓库",
	Long:  `isx fork`,
	Run: func(cmd *cobra.Command, args []string) {
		ForkCmdMain()
	},
}

func ForkCmdMain() {
	fmt.Println("fork")
	currentProject := viper.GetString("current-project.name")
	fmt.Println(currentProject)
	ForkRepository("flink-yun", "TestRepo")
}

func ForkRepository(projectName, newName string) string {
	client := &http.Client{}
	url := common.GithubApiReposDomain + "/isxcode/" + projectName + "/forks"
	req, err := http.NewRequest("POST", url, nil)

	req.Header = common.GitHubHeader(viper.GetString("user.token"))
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败:", err)
		os.Exit(1)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("关闭响应体失败:", err)
		}
	}(resp.Body)

	if resp.StatusCode == http.StatusAccepted {
		fmt.Println("正在处理中，请稍后")
		return "oj8k"
	}
	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("github token权限不足，请重新登录")
		os.Exit(1)
	}
	if resp.StatusCode == http.StatusNotFound {
		fmt.Println("项目不存在")
		os.Exit(1)
	}
	return "omg"
}
