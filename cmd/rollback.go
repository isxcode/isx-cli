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

	// 添加退出选项
	backupDirs = append(backupDirs, "退出")

	// 创建交互式选择器
	prompt := promptui.Select{
		Label:    "请选择要回滚的备份",
		Items:    backupDirs,
		Size:     10,   // 显示最多10个选项
		HideHelp: true, // 隐藏导航提示
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}:",
			Active:   "▶ {{ . | cyan }}",
			Inactive: "  {{ . }}",
			Selected: "✓ {{ . | green }}",
		},
	}

	// 执行选择
	selectedIndex, selectedBackup, err := prompt.Run()
	if err != nil {
		fmt.Printf("选择失败: %v\n", err)
		os.Exit(1)
	}

	// 检查是否选择了退出
	if selectedIndex == len(backupDirs)-1 {
		fmt.Println("已取消操作")
		return
	}

	// 检查当前数据目录是否存在 h2 文件
	if hasH2Files(currentDataPathExpanded) {
		fmt.Println("检测到当前数据目录存在 h2 文件，为了数据安全，不允许回滚覆盖")
		fmt.Println("请先备份当前数据或手动清理 h2 文件后再进行回滚操作")
		return
	}

	// 二次确认
	fmt.Printf("确认要回滚到备份 %s 吗？这将删除当前数据 (y/n): ", selectedBackup)
	var confirm string
	fmt.Scanln(&confirm)

	// 转换为小写进行比较，支持大小写不敏感
	confirm = strings.ToLower(strings.Trim(confirm, " "))

	if confirm != "y" && confirm != "yes" {
		fmt.Println("已中止回滚操作")
		return
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

// hasH2Files 检查指定目录是否存在 h2 数据库文件
func hasH2Files(dirPath string) bool {
	// 检查目录是否存在
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return false
	}

	// 读取目录内容
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return false
	}

	// 检查是否有 .mv.db 或 .trace.db 文件（H2 数据库文件）
	for _, entry := range entries {
		if !entry.IsDir() {
			fileName := entry.Name()
			if strings.HasSuffix(fileName, ".mv.db") ||
				strings.HasSuffix(fileName, ".trace.db") ||
				strings.HasSuffix(fileName, ".lock.db") {
				return true
			}
		}
	}

	return false
}
