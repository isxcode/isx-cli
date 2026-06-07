package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/isxcode/isx-cli/common"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(envCmd)
}

var envCmd = &cobra.Command{
	Use:   "env",
	Short: printCommand("isx env", 40) + "| 安装本地开发环境",
	Long:  `isx env`,
	Run: func(cmd *cobra.Command, args []string) {
		envCmdMain()
	},
}

func envCmdMain() {
	if _, err := exec.LookPath("brew"); err != nil {
		fmt.Println("未检测到 brew，请先安装 Homebrew")
		os.Exit(1)
	}

	installBrewFormula("fnm", "Node.js版本管理工具")
	installBrewFormula("pyenv", "Python版本管理工具")
	installSdkman()
	configureShellEnv()

	fmt.Println("本地开发环境安装完成")
}

func installBrewFormula(formula string, desc string) {
	fmt.Println("开始安装 " + common.WhiteText(formula) + "：" + desc)

	if isBrewFormulaInstalled(formula) {
		fmt.Println(formula + " 已安装")
		return
	}

	runEnvCommand("brew", "install", formula)
}

func isBrewFormulaInstalled(formula string) bool {
	cmd := exec.Command("brew", "list", "--formula", formula)
	err := cmd.Run()
	return err == nil
}

func installSdkman() {
	fmt.Println("开始安装 " + common.WhiteText("sdkman") + "：JDK版本管理工具")

	if isSdkmanInstalled() {
		fmt.Println("sdkman 已安装")
		return
	}

	if isBrewFormulaAvailable("sdkman-cli") {
		runEnvCommand("brew", "install", "sdkman-cli")
		return
	}

	if isBrewFormulaAvailable("sdkman") {
		runEnvCommand("brew", "install", "sdkman")
		return
	}

	fmt.Println("Homebrew 未提供 sdkman 公式，使用 SDKMAN 官方安装脚本")
	runEnvCommand("bash", "-c", "curl -s \"https://get.sdkman.io\" | bash")
}

func isSdkmanInstalled() bool {
	home := common.HomeDir()
	_, err := os.Stat(home + "/.sdkman/bin/sdkman-init.sh")
	return err == nil
}

func isBrewFormulaAvailable(formula string) bool {
	cmd := exec.Command("brew", "info", formula)
	err := cmd.Run()
	return err == nil
}

func runEnvCommand(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("执行失败:", err)
		os.Exit(1)
	}
}

func configureShellEnv() {
	shellName := getCurrentShellName()
	shellConfigPath := getShellConfigPath(shellName)
	shellConfig := buildShellConfig(shellName)

	appendShellConfig(shellConfigPath, shellConfig)

	fmt.Println("环境变量已写入：" + shellConfigPath)
	fmt.Println("当前终端执行以下命令后即可使用 sdk、fnm、pyenv：")
	fmt.Println("source " + shellConfigPath)
}

func getCurrentShellName() string {
	shellPath := os.Getenv("SHELL")
	shellName := filepath.Base(shellPath)
	if shellName == "bash" || shellName == "zsh" {
		return shellName
	}

	if runtime.GOOS == "darwin" {
		return "zsh"
	}
	return "bash"
}

func getShellConfigPath(shellName string) string {
	home := common.HomeDir()
	if shellName == "zsh" {
		return home + "/.zshrc"
	}
	return home + "/.bashrc"
}

func buildShellConfig(shellName string) string {
	return `
# >>> isx env >>>
export SDKMAN_DIR="$HOME/.sdkman"
[[ -s "$SDKMAN_DIR/bin/sdkman-init.sh" ]] && source "$SDKMAN_DIR/bin/sdkman-init.sh"

command -v fnm >/dev/null 2>&1 && eval "$(fnm env --shell ` + shellName + `)"
command -v pyenv >/dev/null 2>&1 && eval "$(pyenv init -)"
# <<< isx env <<<
`
}

func appendShellConfig(path string, shellConfig string) {
	const startMark = "# >>> isx env >>>"

	contentBytes, err := os.ReadFile(path)
	if err == nil && strings.Contains(string(contentBytes), startMark) {
		fmt.Println("环境变量配置已存在：" + path)
		return
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("写入环境变量配置失败:", err)
		os.Exit(1)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("关闭配置文件失败:", err)
		}
	}(file)

	_, err = file.WriteString(shellConfig)
	if err != nil {
		fmt.Println("写入环境变量配置失败:", err)
		os.Exit(1)
	}
}
