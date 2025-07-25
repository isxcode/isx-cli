package cmd

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func init() {
	rootCmd.AddCommand(chooseCmd)
}

var chooseCmd = &cobra.Command{
	Use:   "choose",
	Short: printCommand("isx choose", 40) + "| 切换项目",
	Long:  `从isxcode组织中选择开发项目,isx choose`,
	Run: func(cmd *cobra.Command, args []string) {
		chooseCmdMain()
	},
}

func chooseCmdMain() {
	// 定义项目结构体
	type ProjectConfig struct {
		Name          string `mapstructure:"name"`
		Describe      string `mapstructure:"describe"`
		RepositoryURL string `mapstructure:"repository-url"`
		Dir           string `mapstructure:"dir"`
	}

	// 获取项目列表
	var projectList []ProjectConfig
	err := viper.UnmarshalKey("project-list", &projectList)
	if err != nil {
		fmt.Printf("读取项目列表失败: %v\n", err)
		os.Exit(1)
	}

	// 创建可选择的项目列表（只显示已下载的项目）
	var availableProjects []string
	var availableProjectNames []string

	for _, proj := range projectList {
		// 检查项目是否已下载（通过dir字段判断）
		if proj.Dir != "" {
			// 格式化显示项目信息
			option := fmt.Sprintf("%s : %s",
				printCommand(proj.Name, 12),
				proj.Describe)
			availableProjects = append(availableProjects, option)
			availableProjectNames = append(availableProjectNames, proj.Name)
		}
	}

	// 检查是否有可选择的项目
	if len(availableProjects) == 0 {
		fmt.Println("没有可选择的项目，请先使用 'isx clone' 下载项目代码")
		os.Exit(1)
	}

	// 创建交互式选择器
	prompt := promptui.Select{
		Label: "请选择要切换的项目",
		Items: availableProjects,
		Size:  10, // 显示最多10个选项
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}:",
			Active:   "▶ {{ . | cyan }}",
			Inactive: "  {{ . }}",
			Selected: "✓ {{ . | green }}",
		},
	}

	// 执行选择
	selectedIndex, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("选择失败: %v\n", err)
		os.Exit(1)
	}

	// 设置当前的项目
	projectName := availableProjectNames[selectedIndex]
	fmt.Println("切换到项目：" + projectName)
	viper.Set("now-project", projectName)
	viper.WriteConfig()
}
