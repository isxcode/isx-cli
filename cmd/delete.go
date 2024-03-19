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
	"os/exec"
)

type GithubIssueStatus struct {
	Body  string `json:"body"`
	State string `json:"state"`
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: printCommand("isx delete <issue_number>", 65) + "| 删除远程upstream中分支",
	Long:  `isx delete 123`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			fmt.Println("使用方式不对，请重新输入命令")
			os.Exit(1)
		}

		deleteCmdMain(args[0])
	},
}

func deleteCmdMain(issueNumber string) {

	// 判断issue是否被关闭，未关闭的issue不允许删除
	status := getGithubIssueStatus(issueNumber)
	if status != "closed" {
		fmt.Println("issue未关闭，不允许删除")
		os.Exit(1)
	}

	// 需要删除的分支名
	branchName := "GH-" + issueNumber

	// 删除远程分支
	projectName := viper.GetString("current-project.name")
	projectPath := viper.GetString(projectName+".dir") + "/" + projectName
	deleteUpstreamBranch(projectPath, branchName)

	var subRepository []Repository
	viper.UnmarshalKey(viper.GetString("current-project.name")+".sub-repository", &subRepository)
	for _, repository := range subRepository {
		deleteUpstreamBranch(projectPath+"/"+repository.Name, branchName)
	}
}

func deleteUpstreamBranch(path string, branchName string) {

	executeCommand := "git push upstream -d " + branchName
	deleteBranchCmd := exec.Command("bash", "-c", executeCommand)
	deleteBranchCmd.Stdout = os.Stdout
	deleteBranchCmd.Stderr = os.Stderr
	deleteBranchCmd.Dir = path
	err := deleteBranchCmd.Run()
	if err != nil {
		fmt.Println("远程没有需要删除的分支")
	} else {
		fmt.Println(branchName + "删除成功")
	}
}

func getGithubIssueStatus(issueNumber string) string {

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/repos/isxcode/"+viper.GetString("current-project.name")+"/issues/"+issueNumber, nil)

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

	// 读取响应体内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应体失败:", err)
		os.Exit(1)
	}

	// 解析结果
	if resp.StatusCode == http.StatusOK {
		var content GithubIssueStatus
		json.Unmarshal(body, &content)
		return content.State
	} else if resp.StatusCode == http.StatusNotFound {
		fmt.Println("issue不存在")
		os.Exit(1)
	} else if resp.StatusCode == http.StatusGone {
		fmt.Println("issue已删除,请手动删除分支")
		os.Exit(1)
	} else {
		fmt.Println("无法验证token合法性，登录失败")
		os.Exit(1)
	}
	return ""
}
