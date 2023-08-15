package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
)

func init() {
	rootCmd.AddCommand(prCmd)
}

type GithubTitle struct {
	Title string `json:"title"`
}

var prCmd = &cobra.Command{
	Use:   "pr",
	Short: printCommand("isx pr <issue_number>") + "| 提交代码pr",
	Long:  `快速提交pr，举例：isx pr 123`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			fmt.Println("使用方式不对，请重新输入命令")
			os.Exit(1)
		}

		prCmdMain(args[0])
	},
}

func prCmdMain(issueNumber string) {

	branchName := "GH-" + issueNumber

	// 获取issue的title
	title := getGithubIssueTitle(issueNumber)
	if title == "" {
		fmt.Println("缺陷不存在")
		os.Exit(1)
	}

	// 通过api创建pr
	createPr(branchName+" "+title, branchName)

	var subRepository []Repository
	viper.UnmarshalKey(viper.GetString("current-project.name")+".sub-repository", &subRepository)
	for _ = range subRepository {
		createPr(branchName+" "+title, branchName)
	}
}

func createPr(titleName string, branchName string) {

	headers := http.Header{}
	headers.Set("Accept", "application/vnd.github+json")
	headers.Set("Authorization", "Bearer "+viper.GetString("user.token"))
	headers.Set("X-GitHub-Api-Version", "2022-11-28")

	type ReqJSON struct {
		Title    string `json:"title"`
		Body     string `json:"body"`
		HeadRepo string `json:"head_repo"`
		Head     string `json:"head"`
		Base     string `json:"base"`
	}

	reqJson := ReqJSON{
		Title:    titleName,
		Head:     branchName,
		HeadRepo: viper.GetString("user.account") + "/" + viper.GetString("current-project.name"),
		Base:     branchName,
	}

	client := &http.Client{}
	payload, err := json.Marshal(reqJson)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", "https://api.github.com/repos/isxcode/"+viper.GetString("current-project.name")+"/pulls", bytes.NewBuffer(payload))

	req.Header = headers
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

	// 读取响应体内容
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应体失败:", err)
		os.Exit(1)
	}

	// 解析结果
	if resp.StatusCode == http.StatusCreated {
		fmt.Println(branchName + "提交成功")
	} else if resp.StatusCode == http.StatusNotFound {
		fmt.Println("issue不存在")
	} else if resp.StatusCode == http.StatusUnprocessableEntity {
		fmt.Println("没有提交内容或者重复提交")
	} else {
		fmt.Println("无法验证token合法性，登录失败")
	}
}

func getGithubIssueTitle(issueNumber string) string {

	headers := http.Header{}
	headers.Set("Accept", "application/vnd.github+json")
	headers.Set("Authorization", "Bearer "+viper.GetString("user.token"))
	headers.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/repos/isxcode/"+viper.GetString("current-project.name")+"/issues/"+issueNumber, nil)

	req.Header = headers
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

	// 读取响应体内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应体失败:", err)
		os.Exit(1)
	}

	// 解析结果
	if resp.StatusCode == http.StatusOK {
		var content GithubTitle
		err := json.Unmarshal(body, &content)
		if err != nil {
			fmt.Println("解析 JSON 失败:", err)
		}
		return content.Title
	} else if resp.StatusCode == http.StatusNotFound {
		fmt.Println("issue不存在")
		os.Exit(1)
	} else {
		fmt.Println("无法验证token合法性，登录失败")
		os.Exit(1)
	}

	return ""
}
