package main

import (
	"io/ioutil"
	"log"
)

func readFile(filename string) (string, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("读取文件失败")
		return "", err
	}
	return string(bytes), nil
}
