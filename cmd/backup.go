package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func init() {
	rootCmd.AddCommand(backupCmd)
}

var backupCmd = &cobra.Command{
	Use:   "backup <comment>",
	Short: printCommand("isx backup <comment>", 65) + "| 备份项目数据",
	Long:  `备份项目数据库文件，使用指定的备注作为备份名称`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("使用方式不对，请输入备注信息：isx backup <comment>")
			os.Exit(1)
		}
		backupCmdMain(args[0])
	},
}

func backupCmdMain(comment string) {
	// 获取当前项目名称 - 支持新旧配置格式
	projectName := viper.GetString("now-project")
	if projectName == "" {
		projectName = viper.GetString("current-project.name")
	}

	if projectName == "" {
		fmt.Println("请先使用【isx choose】选择项目")
		os.Exit(1)
	}

	// 根据项目名称确定备份路径
	var sourcePath string
	var backupBasePath string

	switch projectName {
	case "spark-yun":
		sourcePath = "~/.zhiqingyun/h2"
		backupBasePath = "~/.zhiqingyun"
	case "torch-yun":
		sourcePath = "~/.zhishuyun/h2"
		backupBasePath = "~/.zhishuyun"
	default:
		fmt.Printf("项目 %s 暂不支持备份功能\n", projectName)
		os.Exit(1)
	}

	// 展开波浪号路径
	sourcePathExpanded := expandPath(sourcePath)
	backupBasePathExpanded := expandPath(backupBasePath)

	// 检查源路径是否存在
	if _, err := os.Stat(sourcePathExpanded); os.IsNotExist(err) {
		fmt.Printf("数据目录 %s 不存在，无需备份\n", sourcePath)
		os.Exit(1)
	}

	// 生成备份目录名（使用用户备注）
	backupDirName := fmt.Sprintf("h2_%s", comment)
	backupPath := filepath.Join(backupBasePathExpanded, backupDirName)

	// 检查备份目录是否已存在
	if _, err := os.Stat(backupPath); err == nil {
		fmt.Printf("备份 %s 已存在，请使用不同的备注\n", backupDirName)
		os.Exit(1)
	}

	// 确保备份基础目录存在
	if err := os.MkdirAll(backupBasePathExpanded, 0755); err != nil {
		fmt.Printf("创建备份目录失败: %v\n", err)
		os.Exit(1)
	}

	// 执行备份（移动目录而不是复制）
	fmt.Printf("正在备份 %s 到 %s...\n", sourcePath, filepath.Join(backupBasePath, backupDirName))

	moveCmd := exec.Command("mv", sourcePathExpanded, backupPath)
	if err := moveCmd.Run(); err != nil {
		fmt.Printf("备份失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("备份成功！备份文件位于: %s\n", filepath.Join(backupBasePath, backupDirName))
}

// expandPath 展开波浪号路径
func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("获取用户主目录失败: %v\n", err)
			os.Exit(1)
		}
		return filepath.Join(homeDir, path[2:])
	}
	return path
}
