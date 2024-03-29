package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/isxcode/isx-cli/common"
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
	Short: printCommand("isx upgrade", 65) + "| 更新isx-cli",
	Long:  `isx upgrade`,
	Run: func(cmd *cobra.Command, args []string) {
		upgradeCmdMain()
	},
}

type Project struct {
	Name       string `json:"name"`
	Describe   string `json:"describe"`
	Dir        string `json:"dir"`
	Repository struct {
		URL      string `json:"url"`
		Download string `json:"download"`
	} `json:"repository"`
	SubRepository []string `json:"sub-repository"`
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.Contains(s, item) {
			return true
		}
	}
	return false
}

func upgradeCmdMain() {

	// 获取当前版本中的版本号
	oldVersion := viper.GetString("version.number")

	// 判断是否有至匠云模块，没有则直接添加
	projectList := viper.GetStringSlice("project-list")
	if !contains(projectList, "tools-yun") {
		projectList = append(projectList, "tools-yun")
		viper.Set("project-list", projectList)
		toolsYun := Project{
			Name:     "tools-yun",
			Describe: "至匠云，导航中心",
			Dir:      "",
			Repository: struct {
				URL      string `json:"url"`
				Download string `json:"download"`
			}{
				URL:      "https://github.com/isxcode/tools-yun.git",
				Download: "no",
			},
			SubRepository: []string{},
		}
		viper.Set("tools-yun", toolsYun)
		viper.WriteConfig()
	}

	// 获取github中的版本号
	if !contains(projectList, "pytorch-yun") {

		// 在这次更新中,顺带更新项目描述
		viper.Set("flink-yun.describe", "至流云 - 流数据管理平台")
		viper.Set("spark-yun.describe", "至轻云 - 超轻量级大数据平台")
		viper.Set("isx-cli.describe", "至行   - 组织开发专用工具")
		viper.Set("tools-yun.describe", "至匠云 - 企业产品工具导航平台")

		// 项目
		projectList = append(projectList, "pytorch-yun")
		viper.Set("project-list", projectList)
		viper.WriteConfig()

		// 添加配置
		home := common.HomeDir()
		_, err := os.Stat(home + "/.isx/isx-config.yml")
		if !os.IsNotExist(err) {
			pytorchYunStr := "pytorch-yun:\n" +
				"    name: pytorch-yun\n" +
				"    describe: 至慧云 - AI训练部署平台\n" +
				"    dir: \n" +
				"    repository:\n" +
				"        url: https://github.com/isxcode/pytorch-yun.git\n" +
				"        download: no\n" +
				"    sub-repository:\n" +
				"        - url: https://github.com/isxcode/pytorch-yun-vip.git\n" +
				"          name: pytorch-yun-vip"
			file, err := os.OpenFile(home+"/.isx/isx-config.yml", os.O_APPEND|os.O_WRONLY, 0644)
			defer file.Close()
			_, err = file.WriteString(pytorchYunStr)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
		}
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/repos/isxcode/isx-cli/releases/latest", nil)
	if err != nil {
		fmt.Println("创建请求失败:", err)
		os.Exit(1)
	}

	req.Header = common.GitHubHeader(viper.GetString("user.token"))
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
