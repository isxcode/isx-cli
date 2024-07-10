/*
Copyright © 2024 jamie HERE <EMAIL ADDRESS>
*/
package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"io"
	"os"
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

var key = []byte("isxcode")

func Encrypt(token string) string {
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
		fmt.Println("请先登录")
		os.Exit(1)
	}
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
