package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/isxcode/isx-cli/common"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"strconv"
)

func init() {
	rootCmd.AddCommand(issueCmd)
}

var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: printCommand("isx issue", 40) + "| 选择任务",
	Long:  `交互式显示当前仓库分配给您的issue列表，支持光标选择并自动切换到对应分支`,
	Run: func(cmd *cobra.Command, args []string) {
		IssueCmdMain()
	},
}

func IssueCmdMain() {
	username := viper.GetString("user.account")
	// 获取当前项目名称 - 支持新旧配置格式
	currentProject := viper.GetString("now-project")
	if currentProject == "" {
		currentProject = viper.GetString("current-project.name")
	}

	if currentProject == "" {
		fmt.Println("请先使用【isx choose】选择项目")
		os.Exit(1)
	}

	issueList := GetIssueList(currentProject, username)
	if len(issueList) == 0 {
		fmt.Println("当前没有分配给您的issue")
		return
	}

	// 创建选择项列表
	var items []string
	for _, issue := range issueList {
		items = append(items, fmt.Sprintf("GH-%-5d | %s", issue.Number, issue.Title))
	}

	// 创建交互式选择器
	prompt := promptui.Select{
		Label: "请选择要切换的issue",
		Items: items,
		Size:  10, // 显示最多10个选项
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}:",
			Active:   "▶ {{ . | cyan }}",
			Inactive: "  {{ . }}",
			Selected: "✓ {{ . | green }}",
		},
	}

	// 执行选择
	index, selectedItem, err := prompt.Run()
	if err != nil {
		fmt.Printf("选择失败: %v\n", err)
		os.Exit(1)
	}

	// 获取选中的issue号码
	selectedIssue := issueList[index]
	issueNumber := strconv.Itoa(selectedIssue.Number)

	fmt.Printf("已选择issue: %s\n", selectedItem)
	fmt.Printf("正在切换到分支 GH-%s...\n", issueNumber)

	// 调用checkout命令切换到对应分支
	checkoutCmdMain(issueNumber)
}

type IssueListResp struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
}

func GetIssueList(projectName, username string) []IssueListResp {
	client := &http.Client{}
	url := common.GithubApiReposDomain + "/isxcode/" + projectName + "/issues?state=open&assignee=" + username
	req, err := http.NewRequest("GET", url, nil)

	req.Header = common.GitHubHeader(common.GetToken())
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
