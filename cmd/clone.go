package cmd

import (
	"fmt"
	"github.com/isxcode/isx-cli/common"
	"github.com/isxcode/isx-cli/github"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"strings"
)

var projectNumber int
var projectPath string
var projectName string

func init() {
	rootCmd.AddCommand(cloneCmd)
}

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: printCommand("isx clone", 65) + "| 下载项目代码",
	Long:  `isx clone`,
	Run: func(cmd *cobra.Command, args []string) {
		cloneCmdMain()
	},
}

func cloneCmdMain() {

	// 判断用户是否登录
	isLogin := common.CheckUserAccount(common.GetToken())
	if !isLogin {
		fmt.Println("请先登录")
		os.Exit(1)
	}

	// 选择项目编号
	inputProjectNumber()

	// 输入安装路径
	inputProjectPath()

	// 下载项目代码
	cloneProjectCode()

	// 保存配置
	saveConfig()
}

func inputProjectNumber() {
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

	if len(projectList) == 0 {
		fmt.Println("没有可用的项目")
		os.Exit(1)
	}

	// 创建项目选项列表，格式化显示
	var projectOptions []string
	for _, proj := range projectList {
		// 格式化显示项目信息
		option := fmt.Sprintf("%s [%s] : %s",
			printCommand(proj.Name, 14),
			printCommand(proj.RepositoryURL, 45),
			proj.Describe)
		projectOptions = append(projectOptions, option)
	}

	// 创建交互式选择器
	prompt := promptui.Select{
		Label: "请选择要下载的项目",
		Items: projectOptions,
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

	// 设置选中的项目
	projectNumber = selectedIndex
	projectName = projectList[selectedIndex].Name
}

func inputProjectPath() {
	currentWorkDir := common.CurrentWorkDir()

	// 输入安装路径，默认为Y
	fmt.Printf("是否安装在当前路径(%s)下? (Y/n) [默认: Y] ", common.WhiteText(currentWorkDir))
	var flag = ""
	fmt.Scanln(&flag)
	flag = strings.Trim(flag, " ")

	// 如果直接回车，默认为Y
	if flag == "" {
		flag = "Y"
	} else {
		flag = strings.ToUpper(flag)
	}

	// 只接受Y或N
	if flag != "Y" && flag != "N" {
		fmt.Println("输入值异常，请输入 Y 或 N")
		os.Exit(1)
	}

	if flag == "Y" {
		projectPath = currentWorkDir
	} else {
		// flag == "N"
		for {
			fmt.Print("请输入安装路径: ")
			fmt.Scanln(&projectPath)
			projectPath = strings.Trim(projectPath, " ")

			if projectPath == "" {
				fmt.Println("路径不能为空，请重新输入")
				continue
			}

			// 处理路径格式
			projectPath = processPath(projectPath)

			// 检查目录是否存在，如果不存在尝试创建
			if err := ensureDirectoryExists(projectPath); err != nil {
				fmt.Printf("路径处理失败: %v，请重新输入\n", err)
				continue
			}

			break
		}
	}

	// 检查项目是否已存在
	if _, err := os.Stat(projectPath + "/" + projectName); err == nil {
		fmt.Println("项目已存在，请重新选择目录")
		os.Exit(1)
	}
}

// processPath 处理路径格式，支持 ~ 和相对路径
func processPath(path string) string {
	// 支持克隆路径替换～为当前用户目录
	if strings.HasPrefix(path, "~/") {
		path = strings.Replace(path, "~", common.HomeDir(), 1)
	}

	// 处理相对路径
	if !strings.HasPrefix(path, "/") && !strings.Contains(path, ":") {
		// 相对路径，转换为绝对路径
		currentDir := common.CurrentWorkDir()
		path = currentDir + "/" + path
	}

	// 统一路径分隔符
	path = strings.ReplaceAll(path, "\\", "/")

	return path
}

// ensureDirectoryExists 确保目录存在，如果不存在则尝试创建
func ensureDirectoryExists(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// 询问是否创建目录
		fmt.Printf("目录 %s 不存在，是否创建? (Y/n) [默认: Y] ", path)
		var createFlag string
		fmt.Scanln(&createFlag)
		createFlag = strings.Trim(createFlag, " ")

		if createFlag == "" {
			createFlag = "Y"
		} else {
			createFlag = strings.ToUpper(createFlag)
		}

		if createFlag == "Y" {
			err = os.MkdirAll(path, 0755)
			if err != nil {
				return fmt.Errorf("创建目录失败: %v", err)
			}
			fmt.Printf("目录 %s 创建成功\n", path)
		} else {
			return fmt.Errorf("目录不存在且用户选择不创建")
		}
	} else if err != nil {
		return fmt.Errorf("检查目录失败: %v", err)
	}

	return nil
}

func cloneCode(isxcodeRepository string, path string, name string, isMain bool) {

	// 替换下载链接
	isxcodeRepository = strings.Replace(isxcodeRepository, "https://", "https://"+common.GetToken()+"@", -1)

	// 下载主项目代码
	executeCommand := "git clone -b main " + isxcodeRepository
	cloneCmd := exec.Command("bash", "-c", executeCommand)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr
	cloneCmd.Dir = path
	err := cloneCmd.Run()
	if err != nil {
		fmt.Println("执行失败:", err)
		os.Exit(1)
	} else {
		fmt.Println(name + "下载成功")
	}

	// 将origin改为个人的
	userRepository := strings.Replace(isxcodeRepository, "isxcode", viper.GetString("user.account"), -1)
	updateOriginCommand := "git remote set-url origin " + userRepository + " && git fetch origin"
	fmt.Println(updateOriginCommand)
	updateOriginCmd := exec.Command("bash", "-c", updateOriginCommand)
	updateOriginCmd.Stdout = os.Stdout
	updateOriginCmd.Stderr = os.Stderr
	updateOriginCmd.Dir = path + "/" + name
	updateOriginCmd.Run()

	// 添加upstream仓库
	addUpstreamCommand := "git remote add upstream " + isxcodeRepository + " && git fetch upstream"
	fmt.Println(addUpstreamCommand)
	addUpstreamCmd := exec.Command("bash", "-c", addUpstreamCommand)
	addUpstreamCmd.Stdout = os.Stdout
	addUpstreamCmd.Stderr = os.Stderr
	addUpstreamCmd.Dir = path + "/" + name
	addUpstreamCmd.Run()

	// main分支映射到isxcode仓库中
	linkUpstreamCommand := "git branch --set-upstream-to=upstream/main main"
	fmt.Println(linkUpstreamCommand)
	linkUpstreamCmd := exec.Command("bash", "-c", linkUpstreamCommand)
	linkUpstreamCmd.Stdout = os.Stdout
	linkUpstreamCmd.Stderr = os.Stderr
	linkUpstreamCmd.Dir = path + "/" + name
	linkUpstreamCmd.Run()
}

func cloneProjectCode() {
	// 定义项目结构体，包含 sub-repository
	type SubRepository struct {
		Name string `mapstructure:"name"`
		Url  string `mapstructure:"url"`
	}

	type ProjectConfig struct {
		Name          string          `mapstructure:"name"`
		Describe      string          `mapstructure:"describe"`
		RepositoryURL string          `mapstructure:"repository-url"`
		Dir           string          `mapstructure:"dir"`
		SubRepository []SubRepository `mapstructure:"sub-repository"`
	}

	// 获取项目列表
	var projectList []ProjectConfig
	err := viper.UnmarshalKey("project-list", &projectList)
	if err != nil {
		fmt.Printf("读取项目列表失败: %v\n", err)
		os.Exit(1)
	}

	// 找到当前选择的项目
	if projectNumber >= len(projectList) {
		fmt.Println("项目索引超出范围")
		os.Exit(1)
	}

	selectedProject := projectList[projectNumber]
	mainRepository := selectedProject.RepositoryURL

	if mainRepository == "" {
		fmt.Println("项目仓库URL为空")
		os.Exit(1)
	}

	// 下载主项目代码
	if !github.IsRepoForked(viper.GetString("user.account"), projectName) {
		github.ForkRepository("isxcode", projectName, "")
		cloneCode(mainRepository, projectPath, projectName, true)
	} else {
		cloneCode(mainRepository, projectPath, projectName, true)
	}

	// 下载子项目代码 - 支持新配置格式
	var subRepositories []SubRepository

	// 首先尝试从新配置格式获取 sub-repository
	subRepositories = selectedProject.SubRepository

	// 如果新配置格式没有找到，尝试旧配置格式（向后兼容）
	if len(subRepositories) == 0 {
		var legacySubRepository []Repository
		viper.UnmarshalKey(projectName+".sub-repository", &legacySubRepository)
		// 转换旧格式到新格式
		for _, repo := range legacySubRepository {
			subRepositories = append(subRepositories, SubRepository{
				Name: repo.Name,
				Url:  repo.Url,
			})
		}
	}

	// 下载所有子项目
	for _, repository := range subRepositories {
		fmt.Printf("正在处理子项目: %s\n", repository.Name)
		if !github.IsRepoForked(viper.GetString("user.account"), repository.Name) {
			fmt.Printf("Fork子项目: %s\n", repository.Name)
			forkRepository := github.ForkRepository("isxcode", repository.Name, "")
			if forkRepository {
				fmt.Printf("开始下载子项目: %s\n", repository.Name)
				cloneCode(repository.Url, projectPath+"/"+projectName, repository.Name, false)
			} else {
				fmt.Printf("Fork子项目失败: %s\n", repository.Name)
			}
		} else {
			fmt.Printf("子项目已Fork，直接下载: %s\n", repository.Name)
			cloneCode(repository.Url, projectPath+"/"+projectName, repository.Name, false)
		}
	}
}

func saveConfig() {
	// 直接更新指定项目的dir字段，不影响其他字段
	projectList := viper.Get("project-list").([]interface{})

	// 遍历项目列表，找到对应项目并更新dir字段
	for i, item := range projectList {
		if project, ok := item.(map[string]interface{}); ok {
			if name, exists := project["name"]; exists && name == projectName {
				// 保存项目的实际路径（项目根目录）
				project["dir"] = projectPath + "/" + projectName
				projectList[i] = project
				break
			}
		}
	}

	// 保存更新后的项目列表
	viper.Set("project-list", projectList)
	viper.Set("now-project", projectName)
	viper.WriteConfig()
}
