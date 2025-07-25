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

var deleteForceFlag bool
var deleteAllFlag bool

func init() {
	deleteCmd.Flags().BoolVarP(&deleteForceFlag, "force", "f", false, "Force delete")
	deleteCmd.Flags().BoolVarP(&deleteAllFlag, "all", "a", false, "Delete upstream/origin/local")
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: printCommand("isx delete <issue_number>", 40) + "| 删除远程分支",
	Long:  `isx delete 123、isx delete 123 -f 强行删除、 isx delete 123 -a 删除所有和自己相关的需求分支`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			fmt.Println("使用方式不对，请输入：isx delete <issus_number>")
			os.Exit(1)
		}

		deleteCmdMain(args[0])
	},
}

func deleteCmdMain(issueNumber string) {

	// 获取当前项目名称 - 支持新旧配置格式
	projectName := viper.GetString("now-project")
	if projectName == "" {
		projectName = viper.GetString("current-project.name")
	}

	if projectName == "" {
		fmt.Println("请先使用【isx choose】选择项目")
		os.Exit(1)
	}

	// 判断issue是否被关闭，未关闭的issue不允许删除
	if !deleteForceFlag {
		status := getGithubIssueStatus(issueNumber, projectName)
		if status != "closed" {
			fmt.Println("issue未关闭，不允许删除")
			os.Exit(1)
		}
	}

	// 需要删除的分支名
	branchName := "GH-" + issueNumber

	// 获取项目路径 - 支持新旧配置格式
	var projectPath string

	// 尝试新配置格式：从 project-list 数组中查找
	type ProjectConfig struct {
		Name string `mapstructure:"name"`
		Dir  string `mapstructure:"dir"`
	}

	var projectList []ProjectConfig
	configErr := viper.UnmarshalKey("project-list", &projectList)
	if configErr == nil {
		for _, proj := range projectList {
			if proj.Name == projectName {
				projectPath = proj.Dir
				break
			}
		}
	}

	// 如果新配置格式没找到，尝试旧配置格式
	if projectPath == "" {
		projectPath = viper.GetString(projectName + ".dir")
		if projectPath != "" {
			projectPath = projectPath + "/" + projectName
		}
	}

	if projectPath == "" {
		fmt.Printf("项目 %s 未下载，请先使用【isx clone】下载项目代码\n", projectName)
		os.Exit(1)
	}
	deleteUpstreamBranch(projectPath, branchName)

	subRepository := GetSubRepositories(projectName)
	for _, repository := range subRepository {
		deleteUpstreamBranch(projectPath+"/"+repository.Name, branchName)
	}

	// 如果-a的话，删除和自己相关的所有分支
	if deleteAllFlag {
		deleteOriginBranch(projectPath, branchName)
		for _, repository := range subRepository {
			deleteOriginBranch(projectPath+"/"+repository.Name, branchName)
		}
		fmt.Println("请执行 isx git branch -D " + branchName + " 命令删除本地分支 ")
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

func deleteOriginBranch(path string, branchName string) {

	executeCommand := "git push origin -d " + branchName
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

func getGithubIssueStatus(issueNumber string, projectName string) string {

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/repos/isxcode/"+projectName+"/issues/"+issueNumber, nil)

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
