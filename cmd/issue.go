/*
Copyright © 2024 jamie HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/isxcode/isx-cli/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
)

func init() {
	rootCmd.AddCommand(issueCmd)
}

var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: printCommand("isx issue", 65) + "| 列出当前仓库分配给您的issue",
	Long:  `isx issue`,
	Run: func(cmd *cobra.Command, args []string) {
		IssueCmdMain()
	},
}

func IssueCmdMain() {
	username := viper.GetString("user.account")
	currentProject := viper.GetString("current-project.name")
	issueList := GetIssueList(currentProject, username)
	if len(issueList) == 0 {
		fmt.Println("当前没有issue")
	} else {
		for _, issue := range issueList {
			fmt.Printf("💚GH-%-5d | %s \n", issue.Number, issue.Title)
		}
	}
}

type IssueListResp struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
}

func GetIssueList(projectName, username string) []IssueListResp {
	client := &http.Client{}
	url := common.GithubApiReposDomain + "/isxcode/" + projectName + "/issues?state=open&assignee=" + username
	req, err := http.NewRequest("GET", url, nil)

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

	// 解析结果
	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Failed to read response body:", err)
			os.Exit(1)
		}

		var data []IssueListResp
		if err := json.Unmarshal(body, &data); err != nil {
			fmt.Println("Failed to parse JSON response:", err)
			os.Exit(1)
		}
		return data
	} else {
		if resp.StatusCode == http.StatusUnauthorized {
			fmt.Println("github token权限不足，请重新登录")
			os.Exit(1)
		} else {
			fmt.Println("获取issue列表失败")
			fmt.Println("状态码:", resp.StatusCode)
			os.Exit(1)
		}
	}
	return nil
}
