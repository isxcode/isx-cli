package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(postCmd)
}

var postCmd = &cobra.Command{
	Use:   "post <category> <title>",
	Short: printCommand("isx post <category> <title>", 40) + "| 新建博客",
	Long:  `isx post <category> <title>`,
	Run: func(cmd *cobra.Command, args []string) {
		postCmdMain(args)
	},
}

func postCmdMain(args []string) {
	if len(args) < 2 {
		fmt.Println("使用方式不对，请输入：isx post <category> <title>")
		os.Exit(1)
	}

	titleFirst := args[0]
	titleLast := strings.Join(args[1:], " ")
	folder := getPostFolder(titleFirst)

	if folder == "" {
		fmt.Println("该分类不支持")
		os.Exit(1)
	}

	title := titleFirst + " " + titleLast
	postPath := folder + "/" + titleFirst + "/" + title

	cmd := exec.Command("npx", "hexo", "new", titleFirst, "-p", postPath, title)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = getBlogDir()

	fmt.Println("npx hexo new " + titleFirst + " -p \"" + postPath + "\" \"" + title + "\"")

	err := cmd.Run()
	if err != nil {
		fmt.Println("执行失败:", err)
		os.Exit(1)
	}

	fmt.Println("创建博客成功")
}

func getPostFolder(titleFirst string) string {
	postFolders := map[string][]string{
		"github":    []string{"docsify", "git", "github", "hexo", "markdown", "vscode"},
		"hadoop":    []string{"hadoop", "hbase", "hive", "flink", "spark", "kafka", "sqoop", "canal", "zookeeper", "atlas", "cloudera", "solr"},
		"kubernetes": []string{"go", "golang", "kubernetes", "docker", "rancher", "jenkins"},
		"os":        []string{"linux", "mac", "windows", "ngrok", "clash"},
		"pytorch":   []string{"anaconda", "pytorch", "python", "pycharm", "scrapy"},
		"spring":    []string{"java", "spring", "idea", "gradle", "maven", "rabbitmq", "redis"},
		"vue":       []string{"node", "typescript", "vue", "webstorm", "vite", "nginx", "html", "sass", "antdesign", "element"},
		"db":        []string{"mongodb", "mysql", "oracle", "sqlserver", "postgre", "h2", "clickhouse", "doris", "starrocks"},
	}

	for folder, categories := range postFolders {
		for _, category := range categories {
			if titleFirst == category {
				return folder
			}
		}
	}

	return ""
}
