package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"io"
	"os"
	"strings"
)

func HomeDir() string {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return home
}

func CurrentWorkDir() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return dir
}

// GenerateEncryptionKey 生成16字节的随机加密密钥
func GenerateEncryptionKey() string {
	key := make([]byte, 16)
	_, err := rand.Read(key)
	if err != nil {
		fmt.Println("生成加密密钥失败:", err)
		os.Exit(1)
	}
	return hex.EncodeToString(key)
}

// GetEncryptionKey 从配置文件获取加密密钥
func GetEncryptionKey() []byte {
	secret := viper.GetString("user.secret")
	if secret == "" {
		fmt.Println("加密密钥不存在，请重新初始化配置")
		os.Exit(1)
	}
	key, err := hex.DecodeString(secret)
	if err != nil {
		fmt.Println("解析加密密钥失败:", err)
		os.Exit(1)
	}
	return key
}

func Encrypt(token string) string {
	key := GetEncryptionKey()
	ciphertext, err := encryptAES([]byte(token), key)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return ciphertext
}

func GetToken() string {
	token := viper.GetString("user.token")
	if token == "" {
		fmt.Println("未登录，请先使用【isx login】登录账号")
		os.Exit(1)
	}
	if strings.HasPrefix(token, "ghp_") {
		encryptToken := Encrypt(token)
		viper.Set("user.token", encryptToken)
		viper.WriteConfig()
		return token
	}
	key := GetEncryptionKey()
	s, err := decryptAES(token, key)
	if err != nil {
		fmt.Println("解密失败...", err)
		os.Exit(1)
	}
	return s
}

func encryptAES(plaintext, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return hex.EncodeToString(ciphertext), nil
}

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

func Parse(reader io.Reader, v any) {
	body, err := io.ReadAll(reader)
	if err != nil {
		fmt.Println("读取响应体失败:", err)
		os.Exit(1)
	}
	err = json.Unmarshal(body, v)
	if err != nil {
		fmt.Println("解析 JSON 失败:", err)
		os.Exit(1)
	}
}

func ToJsonString(v any) string {
	return string(ToJsonBytes(v))
}

func ToJsonBytes(v any) []byte {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		fmt.Println("解析 JSON 失败:", err)
		os.Exit(1)
	}
	return jsonBytes
}
