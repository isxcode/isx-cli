package git

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/exec"
	"strings"
)

const BranchTemplate = "当前项目: %s\n当前分支: %s\n项目路径: %s\n"

func GetCurrentBranchName(projectName string, abortOnFailure bool) string {
	projectPath := viper.GetString(projectName+".dir") + "/" + projectName

	executeCommand := "git symbolic-ref --short HEAD"
	branchCmd := exec.Command("bash", "-c", executeCommand)
	branchCmd.Dir = projectPath
	output, err := branchCmd.Output()
	if err != nil {
		log.Println("获取当前项目分支名失败:", err)
		if abortOnFailure {
			fmt.Printf(BranchTemplate, projectName, "获取分支名称失败", projectPath)
			os.Exit(1)
		}
		return "获取分支名称失败"
	}
	return strings.Split(string(output), "\n")[0]
}

func PushBranchToOrigin(branchName, path string) {
	pushBranch(branchName, "origin", path)
}

func PushBranchToUpstream(branchName, path string) {
	pushBranch(branchName, "upstream", path)
}

func pushBranch(branchName, repositoryAlias, path string) {
	pushOriginCmd := exec.Command("bash", "-c", "git push "+repositoryAlias+" "+branchName)
	pushOriginCmd.Stdout = os.Stdout
	pushOriginCmd.Stderr = os.Stderr
	pushOriginCmd.Dir = path
	err := pushOriginCmd.Run()
	if err != nil {
		fmt.Println("执行失败:", err)
		os.Exit(1)
	} else {
		fmt.Println(branchName + "，已推到" + repositoryAlias + "仓库")
	}
}
