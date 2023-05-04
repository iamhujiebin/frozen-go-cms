package main

import (
	"fmt"
	"frozen-go-cms/route"
	"net/http"
)

const (
	PORT        = 8086
	STATIC_PORT = 8084
)

func main() {
	go static()                     // 静态服务器
	r := route.InitRouter()         // 注册路由
	r.Run(fmt.Sprintf(":%d", PORT)) // 启动服务
}

func static() {
	fs := http.FileServer(http.Dir("build/"))
	http.Handle("/", fs)

	http.ListenAndServe(fmt.Sprintf(":%d", STATIC_PORT), nil)
}
