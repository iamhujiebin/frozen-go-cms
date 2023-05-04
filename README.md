# CMS系统-服务端

1. 提供静态文件服务器
2. 提供图片文件服务器
3. 提供用户系统/文章系统等接口

# 初始化步骤

1. go mod init frozen-go-cms
2. go get 报错的
3. 暂时还是先依赖hilo-common

# 项目架构

+ .env 项目环境配置,正服/测服会用服务器上的配置覆盖
+ run.sh 守护脚本
+ Makefile 构建文件
+ main.go 入口文件
+ *.ini 配置文件,测服用debug.ini,正服用release.ini
+ 目录分层说明(`着重注意命名规则`,下划线_分割职责)
    + test: 单元测试目录
    + _const: 常量
        + 子目录 enum,放业务定义常量
            + 命名规则 `*_e`,如game_e
        + 子目录 redis_key,放redis的key
            + 命名规则 `*_k`,如game_k
    + cron: 定时任务
        + 统一入口cron.go
        + 子目录命名规则 `*_cron`,如game_cron
    + mycontext: 上下文
    + myerr: 错误码
        + 子目录 bizerr,放业务错误
    + mylogrus: 日志包
    + req: 请求参数
        + request.go: 定义一些常用方法
        + 子目录jwt: 鉴权相关
        + 子目录
            + 命名规则 `*_req`,如game_req
    + resp: 返回结构体,通用结构体
    + cv: 客户端需要结构体
        + 子目录
            + 命名规则 `*_cv`,如game_cv
    + route: 路由定义
        + 根目录
            + errorHandler.go 错误处理
            + middleHandle.go 中间件
            + router.go 路由定义
            + util.go 工具包
        + 子目录,业务路由定义
            + 命名规则 `*_r`,如game_r
    + resource: 资源层
        + config: 配置相关
        + consul: 注册中心
        + mysql: 数据库
        + redisCli: 缓存
    + domain: 领域层
        + ctx.go: 定义ctxAndDb
        + model.go: 通用model
        + event.go: 抽象定义event
        + cache: 缓存层
            + moduleRedis.go 带有model的通用redis方法
            + 子目录
                + 命名规则 `*_c`,如user_c
        + event: 事件层
            + base.go: 定义事件base
            + 子目录
                + 命名规则 `*_ev`,如game_ev
        + model: 模型层
            + repo.go 持久化
            + 子目录
                + 命名规则 `*_m`,如game_m/user_m
        + service: 服务层,主要是开启事务和发事件
            + service.go: 事务/事件封装
            + 子目录
                + 命名规则 `*_s`,如game_s