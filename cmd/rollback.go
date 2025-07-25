package cmd

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func init() {
	rootCmd.AddCommand(rollbackCmd)
}

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: printCommand("isx rollback", 40) + "| 回滚项目资源",
	Long:  `回滚项目数据库文件到指定备份`,
	Run: func(cmd *cobra.Command, args []string) {
		rollbackCmdMain()
	},
}

func rollbackCmdMain() {
	// 获取当前项目名称 - 支持新旧配置格式
	projectName := viper.GetString("now-project")
	if projectName == "" {
		projectName = viper.GetString("current-project.name")
	}

	if projectName == "" {
		fmt.Println("请先使用【isx choose】选择项目")
		os.Exit(1)
	}

	// 根据项目名称确定路径
	var currentDataPath string
	var backupBasePath string

	switch projectName {
	case "spark-yun":
		currentDataPath = "~/.zhiqingyun/h2"
		backupBasePath = "~/.zhiqingyun"
	case "torch-yun":
		currentDataPath = "~/.zhishuyun/h2"
		backupBasePath = "~/.zhishuyun"
	default:
		fmt.Printf("项目 %s 暂不支持回滚功能\n", projectName)
		os.Exit(1)
	}

	// 展开波浪号路径
	currentDataPathExpanded := expandPath(currentDataPath)
	backupBasePathExpanded := expandPath(backupBasePath)

	// 检查备份基础目录是否存在
	if _, err := os.Stat(backupBasePathExpanded); os.IsNotExist(err) {
		fmt.Printf("备份目录 %s 不存在\n", backupBasePath)
		os.Exit(1)
	}

	// 获取所有备份目录
	backupDirs, err := getBackupDirectories(backupBasePathExpanded)
	if err != nil {
		fmt.Printf("获取备份目录列表失败: %v\n", err)
		os.Exit(1)
	}

	if len(backupDirs) == 0 {
		fmt.Println("没有找到任何备份文件")
		os.Exit(1)
	}

	// 创建交互式选择器
	prompt := promptui.Select{
		Label: "请选择要回滚的备份",
		Items: backupDirs,
		Size:  10, // 显示最多10个选项
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}:",
			Active:   "▶ {{ . | cyan }}",
			Inactive: "  {{ . }}",
			Selected: "✓ {{ . | green }}",
		},
	}

	// 执行选择
	_, selectedBackup, err := prompt.Run()
	if err != nil {
		fmt.Printf("选择失败: %v\n", err)
		os.Exit(1)
	}

	// 二次确认
	confirmPrompt := promptui.Prompt{
		Label:     fmt.Sprintf("确认要回滚到备份 %s 吗？这将删除当前数据", selectedBackup),
		IsConfirm: true,
	}

	_, err = confirmPrompt.Run()
	if err != nil {
		fmt.Println("已中止回滚操作")
		os.Exit(1)
	}

	// 执行回滚
	selectedBackupPath := filepath.Join(backupBasePathExpanded, selectedBackup)

	fmt.Printf("正在回滚到备份: %s\n", selectedBackup)

	// 1. 删除当前数据目录（如果存在）
	if _, err := os.Stat(currentDataPathExpanded); err == nil {
		fmt.Println("删除当前数据目录...")
		if err := os.RemoveAll(currentDataPathExpanded); err != nil {
			fmt.Printf("删除当前数据目录失败: %v\n", err)
			os.Exit(1)
		}
	}

	// 2. 移动备份到当前数据目录（不保留备份）
	fmt.Println("恢复备份数据...")
	moveCmd := exec.Command("mv", selectedBackupPath, currentDataPathExpanded)
	if err := moveCmd.Run(); err != nil {
		fmt.Printf("恢复备份失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("回滚成功！已恢复到备份: %s\n", selectedBackup)
}

// getBackupDirectories 获取所有备份目录
func getBackupDirectories(basePath string) ([]string, error) {
	entries, err := os.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	var backupDirs []string
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "h2_") {
			backupDirs = append(backupDirs, entry.Name())
		}
	}

	// 按时间戳排序（最新的在前）
	sort.Sort(sort.Reverse(sort.StringSlice(backupDirs)))

	return backupDirs, nil
}
