package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/tidwall/gjson"
)

var (
	conf          string
	errcode       map[string]gjson.Result
	err           error
	suburl        string
	dirviewerpath string
	client        *oss.Client
)

func main() {
	conf, err = readFile("config.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	errcode = gjson.Get(conf, "errcode").Map()
	fmt.Println("OSS Go SDK Version: ", oss.Version)

	endpoint := gjson.Get(conf, "endpoint").String()
	accesskeyid := gjson.Get(conf, "accesskeyid").String()
	accesskeysecret := gjson.Get(conf, "accesskeysecret").String()
	linkTimeOut := gjson.Get(conf, "timeOut.link").Int()
	ioTimeOut := gjson.Get(conf, "timeOut.io").Int()

	client, err = oss.New(endpoint, accesskeyid, accesskeysecret, oss.Timeout(linkTimeOut, ioTimeOut), oss.EnableCRC(false))
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	setupCloseHandler()
	suburl = gjson.Get(conf, "suburl").String()
	dirviewerpath = gjson.Get(conf, "dirviewer").String()
	listenandserve := gjson.Get(conf, "listenandserve")

	fmt.Printf("启动 HTTP 服务（端口 %s ）... \n", listenandserve)
	fmt.Println(suburl + dirviewerpath)

	fmt.Println("准备就绪。")
	http.HandleFunc("/", mainHandleFunc)
	http.HandleFunc(suburl+dirviewerpath, dirviewerHandleFunc)
	err = http.ListenAndServe(":"+listenandserve.String(), nil) //设置监听的端口
	if err != nil {
		fmt.Println(err)
	}
}

func mainHandleFunc(w http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Method + " :" + req.Header.Get("X-Forwarded-For") + " " + req.RemoteAddr + " -> " + req.RequestURI)
	fmt.Println("404")
}
func setupCloseHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("收到中止请求，正在退出 ... ")
		fmt.Println("退出。")
		os.Exit(0)
	}()
}
