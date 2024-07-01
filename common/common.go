/*
Copyright © 2024 jamie HERE <EMAIL ADDRESS>
*/
package common

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
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
