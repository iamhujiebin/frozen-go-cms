package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"frozen-go-cms/cron"
	"frozen-go-cms/route"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

const (
	PORT        = 7000
	STATIC_PORT = 7001
	FILE_PORT   = 7002
	DOMAIN      = "http://43.128.31.213"
	UPLOAD_PATH = "uploads/file"
)

func main() {
	cron.Init()
	go static()                         // 静态服务器
	go file()                           // 文件服务器
	r := route.InitRouter()             // 注册路由
	_ = r.Run(fmt.Sprintf(":%d", PORT)) // 启动服务
}

// 静态服务器
func static() {
	var router = gin.Default()
	router.StaticFS("/", http.Dir("build/"))
	_ = router.Run(fmt.Sprintf(":%d", STATIC_PORT))
}

// 文件服务器
func file() {
	var router = gin.Default()
	router.Use(route.Cors())
	router.StaticFS(UPLOAD_PATH, http.Dir(UPLOAD_PATH))
	group := router.Group("file") // 跨域
	group.POST("/upload", uploadFunc)
	group.POST("/upload2", uploadHandler)
	_ = router.Run(fmt.Sprintf(":%d", FILE_PORT))
}

func uploadFunc(c *gin.Context) {
	prefix := DOMAIN + fmt.Sprintf(":%d/", FILE_PORT)
	_, header, err := c.Request.FormFile("file")
	code, message := 0, "success"
	if err != nil {
		code = 1001
		message = err.Error()
	}
	filepath, filename, err := UploadFile(header)
	if err != nil {
		code = 1001
		message = err.Error()
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    code,
		"message": message,
		"data": map[string]interface{}{
			"filepath": filepath,
			"filename": filename,
			"url":      prefix + filepath,
		},
	})
}

func UploadFile(uploadFile *multipart.FileHeader) (filepath, filename string, err error) {
	// 读取文件后缀
	ext := path.Ext(uploadFile.Filename)
	// 读取文件名并md5加密
	name := strings.TrimSuffix(uploadFile.Filename, ext)
	name = MD5V([]byte(name))
	// 拼接新文件名(用时间戳)
	filename = name + "_" + time.Now().Format("20060102150405") + ext
	// 创建路径
	err = os.MkdirAll(UPLOAD_PATH, os.ModePerm)
	if err != nil {
		return
	}
	// 拼接路径和文件名
	filepath = UPLOAD_PATH + "/" + filename

	// 读取上传的文件
	f, err := uploadFile.Open()
	if err != nil {
		return
	}
	defer f.Close()

	// 创建文件
	out, err := os.Create(filepath)
	if err != nil {
		return
	}
	defer out.Close()

	// 拷贝文件
	_, err = io.Copy(out, f)
	if err != nil {
		return
	}
	return
}

// md5加密
func MD5V(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}

const (
	ChunkSize = 1024 * 1024 // 1MB
)

func uploadHandler(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	code, message := 0, "success"
	if err != nil {
		code = 1001
		message = err.Error()
	}
	defer file.Close()

	// 创建目标文件
	targetFile, err := os.Create("./uploads/file2/" + header.Filename)
	if err != nil {
		code = 1001
		message = err.Error()
	}
	defer targetFile.Close()

	// 分片上传
	buf := make([]byte, ChunkSize)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			code = 1002
			message = err.Error()
			break
		}

		if n == 0 {
			break
		}

		_, err = targetFile.Write(buf[:n])
		if err != nil {
			code = 1003
			message = err.Error()
			break
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    code,
		"message": message,
		"data": map[string]interface{}{
			"filepath": "todo",
			"filename": "dodo",
			"url":      "todo",
		},
	})
}
