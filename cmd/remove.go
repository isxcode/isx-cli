package cmd

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var deleteProjectNumber int

func init() {
	rootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: printCommand("isx remove", 40) + "| 删除本地项目",
	Long:  `isx remove`,
	Run: func(cmd *cobra.Command, args []string) {
		removeCmdMain()
	},
}

func removeCmdMain() {

	// 选择项目编号
	inputRemoveProjectNumber()

	// 删除项目
	removeProject()
}

func inputRemoveProjectNumber() {
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

	// 创建可删除的项目列表（只显示已下载的项目）
	var removableProjects []string
	var removableProjectIndices []int

	for i, proj := range projectList {
		// 检查项目是否已下载（通过dir字段判断）
		if proj.Dir != "" {
			// 格式化显示项目信息
			option := fmt.Sprintf("%s [%s] : %s",
				printCommand(proj.Name, 14),
				printCommand(proj.RepositoryURL, 40),
				proj.Describe)
			removableProjects = append(removableProjects, option)
			removableProjectIndices = append(removableProjectIndices, i)
		}
	}

	// 检查是否有可删除的项目
	if len(removableProjects) == 0 {
		fmt.Println("没有可删除的项目，请先使用 【isx clone】 下载项目代码")
		os.Exit(1)
	}

	// 创建交互式选择器
	prompt := promptui.Select{
		Label:    "请选择要删除的项目",
		Items:    removableProjects,
		Size:     10,   // 显示最多10个选项
		HideHelp: true, // 隐藏导航提示
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}:",
			Active:   "▶ {{ . | red }}",
			Inactive: "  {{ . }}",
			Selected: "✗ {{ . | red }}",
		},
	}

	// 执行选择
	selectedIndex, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("选择失败: %v\n", err)
		os.Exit(1)
	}

	// 设置要删除的项目索引
	deleteProjectNumber = removableProjectIndices[selectedIndex]

	// 二次确认
	confirmPrompt := promptui.Prompt{
		Label:     "确认要删除该项目吗",
		IsConfirm: true,
	}

	_, err = confirmPrompt.Run()
	if err != nil {
		fmt.Println("已中止删除操作")
		os.Exit(1)
	}
}

func removeProject() {
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

	// 获取项目目录
	selectedProject := projectList[deleteProjectNumber]
	projectName := selectedProject.Name
	projectFullPath := selectedProject.Dir // 这已经是完整路径，包含项目名

	// 三次确认删除
	deleteProject := "N"
	fmt.Print("确认要删除【" + projectFullPath + "】路径吗?(Y/N) default is N: ")
	fmt.Scanln(&deleteProject)
	if deleteProject == "N" || deleteProject == "n" {
		fmt.Println("已中止")
		os.Exit(1)
	}

	// 更新平台替换projectPath
	removeCommand := ""
	if runtime.GOOS == "windows" {
		projectFullPath = strings.ReplaceAll(projectFullPath, "C:", "/c")
		projectFullPath = strings.ReplaceAll(projectFullPath, " ", "\\ ")
		removeCommand = "rm -rf " + projectFullPath
	} else {
		removeCommand = "rm -rf " + projectFullPath
	}

	removeCmd := exec.Command("bash", "-c", removeCommand)
	removeCmd.Stdout = os.Stdout
	removeCmd.Stderr = os.Stderr
	err = removeCmd.Run()
	if err != nil {
		fmt.Println("执行失败:", err)
		os.Exit(1)
	} else {
		fmt.Println(projectFullPath + "路径已删除")
	}

	// 保存配置：清空project-list中对应项目的dir字段
	projectListRaw := viper.Get("project-list").([]interface{})

	// 遍历项目列表，找到对应项目并清空dir字段
	for i, item := range projectListRaw {
		if project, ok := item.(map[string]interface{}); ok {
			if name, exists := project["name"]; exists && name == projectName {
				// 只清空dir字段
				project["dir"] = ""
				projectListRaw[i] = project
				break
			}
		}
	}

	// 更新配置
	viper.Set("project-list", projectListRaw)
	if viper.GetString("now-project") == projectName {
		viper.Set("now-project", "")
	}
	viper.WriteConfig()
}
