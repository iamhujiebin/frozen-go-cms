package main

import (
	"fmt"
	"frozen-go-cms/route"
)

const (
	PORT = 8086
)

func main() {
	// 静态文件服务器

	r := route.InitRouter()         // 注册路由
	r.Run(fmt.Sprintf(":%d", PORT)) // 启动服务
}
