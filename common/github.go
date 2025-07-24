package common

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const GithubDomain = "https://github.com"
const GithubRawDomain = "https://raw.github.com"
const GithubApiDomain = "https://api.github.com"
const GithubApiReposDomain = "https://api.github.com/repos"
const IsxcodeGithubApiReposDomain = GithubApiReposDomain + "/isxcode"

func GitHubHeader(accessToken string) http.Header {
	headers := http.Header{}
	headers.Set("Accept", "application/vnd.github+json")
	headers.Set("Authorization", "Bearer "+accessToken)
	headers.Set("X-GitHub-Api-Version", "2022-11-28")
	return headers
}

// GitHubUserInfo GitHub用户信息结构体
type GitHubUserInfo struct {
	Login string `json:"login"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func CheckUserAccount(accessToken string) bool {

	client := &http.Client{}
	req, err := http.NewRequest("GET", GithubApiDomain+"/octocat", nil)

	req.Header = GitHubHeader(accessToken)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败:", err)
		os.Exit(1)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// 解析结果
	if resp.StatusCode == http.StatusOK {
		return true
	} else {
		return false
	}
}

// GetGitHubUserInfo 获取GitHub用户信息
func GetGitHubUserInfo(accessToken string) (*GitHubUserInfo, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", GithubApiDomain+"/user", nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header = GitHubHeader(accessToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("关闭响应体失败:", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取用户信息失败，状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %v", err)
	}

	var userInfo GitHubUserInfo
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	return &userInfo, nil
}
