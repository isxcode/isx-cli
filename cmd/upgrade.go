package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
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
	Short: printCommand("isx upgrade", 40) + "| 升级脚手架",
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
	latestVersion := getLatestVersion()
	if latestVersion == "" {
		fmt.Println("获取最新版本失败")
		os.Exit(1)
	}

	// 每次升级都直接重新下载安装
	// 执行更新命令
	executeCommand := "sh -c \"$(curl -fsSL https://isxcode.oss-cn-shanghai.aliyuncs.com/zhixingyun/install.sh)\""
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
	viper.Set("version", latestVersion)
	viper.WriteConfig()
}

// getLatestVersion 获取最新版本，优先使用匿名请求，失败时尝试认证请求
func getLatestVersion() string {
	// 首先尝试匿名请求
	version := tryAnonymousRequest()
	if version != "" {
		return version
	}

	// 匿名请求失败，尝试认证请求
	return tryAuthenticatedRequest()
}

// tryAnonymousRequest 尝试匿名请求获取最新版本
func tryAnonymousRequest() string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/repos/isxcode/isx-cli/releases/latest", nil)
	if err != nil {
		fmt.Println("创建匿名请求失败:", err)
		return ""
	}

	// 设置基本的请求头，但不包含认证信息
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("匿名请求失败:", err)
		return ""
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// 忽略关闭错误
		}
	}(resp.Body)

	return parseVersionResponse(resp, "匿名")
}

// tryAuthenticatedRequest 尝试认证请求获取最新版本
func tryAuthenticatedRequest() string {
	// 检查是否有可用的 token
	token := viper.GetString("user.token")
	if token == "" {
		fmt.Println("匿名请求达到速率限制，且未登录。建议使用 'isx login' 登录以获得更高的速率限制")
		return ""
	}

	// 尝试获取并解密 token
	var actualToken string
	if strings.HasPrefix(token, "ghp_") {
		actualToken = token
	} else {
		// 检查是否有加密密钥
		secret := viper.GetString("user.secret")
		if secret == "" {
			fmt.Println("配置文件损坏，请重新登录")
			return ""
		}

		key, err := hex.DecodeString(secret)
		if err != nil {
			fmt.Println("解析加密密钥失败，请重新登录")
			return ""
		}

		decryptedToken, err := decryptAES(token, key)
		if err != nil {
			fmt.Println("解密 token 失败，请重新登录")
			return ""
		}
		actualToken = decryptedToken
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/repos/isxcode/isx-cli/releases/latest", nil)
	if err != nil {
		fmt.Println("创建认证请求失败:", err)
		return ""
	}

	req.Header = common.GitHubHeader(actualToken)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("认证请求失败:", err)
		return ""
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// 忽略关闭错误
		}
	}(resp.Body)

	return parseVersionResponse(resp, "认证")
}

// parseVersionResponse 解析版本响应
func parseVersionResponse(resp *http.Response, requestType string) string {
	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("%s请求：读取响应体失败: %v\n", requestType, err)
			return ""
		}

		var data map[string]interface{}
		if err := json.Unmarshal(body, &data); err != nil {
			fmt.Printf("%s请求：解析 JSON 失败: %v\n", requestType, err)
			return ""
		}

		if name, ok := data["name"].(string); ok {
			return strings.ReplaceAll(name, "v", "")
		} else {
			fmt.Printf("%s请求：响应中缺少版本信息\n", requestType)
			return ""
		}
	} else if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests {
		if requestType == "匿名" {
			fmt.Printf("%s请求达到速率限制 (状态码: %d)，尝试使用认证请求...\n", requestType, resp.StatusCode)
		} else {
			fmt.Printf("%s请求达到速率限制 (状态码: %d)\n", requestType, resp.StatusCode)
		}
		return ""
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Printf("%s请求：认证失败，请重新登录\n", requestType)
		return ""
	} else {
		fmt.Printf("%s请求失败，状态码: %d\n", requestType, resp.StatusCode)
		return ""
	}
}

// decryptAES 解密 AES 加密的字符串（从 common 包复制，避免循环依赖）
func decryptAES(ciphertext string, key []byte) (string, error) {
	ciphertextBytes, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertextBytes) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertextBytes, ciphertextBytes)

	return string(ciphertextBytes), nil
}
