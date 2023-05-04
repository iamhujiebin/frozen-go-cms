package main

import (
	"fmt"
	"frozen-go-cms/route"
)

const (
	PORT = 8086
)

func main() {
	//cron.Init()                     // 开启定时任务
	r := route.InitRouter()         // 注册路由
	r.Run(fmt.Sprintf(":%d", PORT)) // 启动服务
}
