package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func init() {
	rootCmd.AddCommand(upgradeCmd)
}

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: printCommand("isx upgrade") + "| 更新isx-cli",
	Long:  `isx upgrade`,
	Run: func(cmd *cobra.Command, args []string) {
		upgradeCmdMain()
	},
}

func upgradeCmdMain() {

	// 获取当前版本中的版本号
	oldVersion := viper.GetString("version.number")

	// 获取github中的版本号
	headers := http.Header{}
	headers.Set("Accept", "application/vnd.github+json")
	headers.Set("Authorization", "Bearer "+viper.GetString("user.token"))
	headers.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/repos/isxcode/isx-cli/releases/latest", nil)
	if err != nil {
		fmt.Println("创建请求失败:", err)
		os.Exit(1)
	}

	req.Header = headers
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

	latestVersion := ""

	// 解析结果
	if resp.StatusCode == http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			fmt.Println("github token权限不足，请重新登录")
			os.Exit(1)
		} else {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Failed to read response body:", err)
				os.Exit(1)
			}

			var data map[string]interface{}
			if err := json.Unmarshal(body, &data); err != nil {
				fmt.Println("Failed to parse JSON response:", err)
				os.Exit(1)
			}

			latestVersion = strings.ReplaceAll(data["name"].(string), "v", "")
		}
	} else {
		fmt.Println("获取最新版本失败")
		os.Exit(1)
	}

	// 版本号进行对比
	if oldVersion < latestVersion {

		// 执行更新命令
		executeCommand := "sh -c \"$(curl -fsSL https://raw.githubusercontent.com/isxcode/isx-cli/main/install.sh)\""
		result := exec.Command("bash", "-c", executeCommand)
		result.Stdout = os.Stdout
		result.Stderr = os.Stderr

		err := result.Run()
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("已更新到最新版本：" + latestVersion)
		}

		// 更新配置中的版本信息
		viper.Set("version.number", latestVersion)
		viper.WriteConfig()
	} else {
		fmt.Println("已经是最新版本")
		os.Exit(1)
	}
}
