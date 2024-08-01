package common

import (
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
