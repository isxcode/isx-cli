package github

import (
	"bytes"
	"fmt"
	"github.com/isxcode/isx-cli/common"
	"io"
	"net/http"
	"strings"
)

const createForkUrl = common.GithubApiReposDomain + "/%s/%s/forks"
const isForkedUrl = common.GithubApiReposDomain + "/%s/%s"

func IsRepoForked(account, projectName string) bool {
	resp := Get(fmt.Sprintf(isForkedUrl, account, projectName), nil)
	defer CloseRespBody(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var repo Repository
		common.Parse(resp.Body, &repo)
		return repo.Fork
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return false
	}
	if resp.StatusCode == http.StatusNotFound {
		return false
	}
	return false
}

func ForkRepository(owner, projectName, newName string) bool {
	var reqBody io.Reader = nil
	if len(newName) > 0 {
		newName = strings.Trim(newName, " ")
		if len(newName) > 0 {
			jsonBytes := common.ToJsonBytes(map[string]string{"name": newName})
			reqBody = bytes.NewBuffer(jsonBytes)
		}
	}

	resp := Post(fmt.Sprintf(createForkUrl, owner, projectName), reqBody)
	defer CloseRespBody(resp.Body)

	if resp.StatusCode == http.StatusAccepted {
		return true
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return false
	}
	if resp.StatusCode == http.StatusNotFound {
		return false
	}
	return false
}
