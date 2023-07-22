package ai_r

import (
	"encoding/json"
	"errors"
	"fmt"
	"frozen-go-cms/domain/model/ai_m"
	"frozen-go-cms/hilo-common/domain"
	"frozen-go-cms/hilo-common/mycontext"
	"frozen-go-cms/hilo-common/mylogrus"
	"frozen-go-cms/hilo-common/resource/mysql"
	"frozen-go-cms/resp"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// @Tags AI
// @Summary 获取所有提示词
// @Param Authorization header string true "token"
// @Success 200 {object} []ai_m.AiPromptData
// @Router /v1_0/ai/prompts [get]
func Prompts(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	res := ai_m.GetAllPrompts(model)
	resp.ResponseOk(c, res)
	return myCtx, nil
}

// @Tags AI
// @Summary 获取所有ai生成的图片
// @Param Authorization header string true "token"
// @Success 200 {object} []ai_m.AiImage
// @Router /v1_0/ai/images [get]
func Images(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	res := ai_m.GetAllImages(model)
	resp.ResponseOk(c, res)
	return myCtx, nil
}

type GenImagesReq struct {
	Prompts string `form:"prompts" binding:"required"`
}

// @Tags AI
// @Summary 获取所有ai生成的图片
// @Param Authorization header string true "token"
// @Success 200 {object} ai_m.AiImage
// @Router /v1_0/ai/images [post]
func GenImages(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	var param GenImagesReq
	if err := c.ShouldBindJSON(&param); err != nil {
		return myCtx, err
	}
	model := domain.CreateModelContext(myCtx)
	genResp, payload, err := genImage(param.Prompts)
	if err != nil {
		return myCtx, err
	}
	status := 1
	if len(genResp.Data) <= 0 {
		status = 0
	}
	var images []string
	for _, v := range genResp.Data {
		images = append(images, v.URL)
	}
	var image1, image2 string
	for i := range genResp.Data {
		if i == 0 {
			image1 = genResp.Data[i].URL
		}
		if i == 1 {
			image2 = genResp.Data[i].URL
		}
	}
	response := &ai_m.AiImage{
		Prompt:  param.Prompts,
		Image1:  image1,
		Image2:  image2,
		Status:  status,
		Like:    0,
		Payload: payload,
	}
	if err := ai_m.AddAiImage(model, response); err != nil {
		return myCtx, err
	}
	// 有效的图片才下载
	if status == 1 {
		go persistentImages(images, response.ID)
		resp.ResponseOk(c, response)
	} else {
		return myCtx, errors.New("gen fail,may hit sensitive word")
	}
	return myCtx, nil
}

type genImageReq struct {
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}

type GenImageResp struct {
	Created int `json:"created"`
	Data    []struct {
		URL string `json:"url"`
	} `json:"data"`
}

func genImage(prompt string) (*GenImageResp, string, error) {
	// todo for coding
	//r := &GenImageResp{Created: 0}
	//r.Data = append(r.Data, struct {
	//	URL string `json:"url"`
	//}{URL: "https://oaidalleapiprodscus.blob.core.windows.net/private/org-lx7MeI8vgrXk2dZIpjAehwph/user-p1Gk7zjjfNg1aoVI1pCcZhUQ/img-KMbFixNoky5IEyLVm869VII8.png?st=2023-07-22T14%3A33%3A02Z&se=2023-07-22T16%3A33%3A02Z&sp=r&sv=2021-08-06&sr=b&rscd=inline&rsct=image/png&skoid=6aaadede-4fb3-4698-a8f6-684d7786b067&sktid=a48cca56-e6da-484e-a814-9c849652bcb3&skt=2023-07-21T23%3A28%3A09Z&ske=2023-07-22T23%3A28%3A09Z&sks=b&skv=2021-08-06&sig=MDnW/Dr10GAIUHGaa/P/yedZUjqmzPgSYe7kxuTz39s%3D"})
	//r.Data = append(r.Data, struct {
	//	URL string `json:"url"`
	//}{URL: "https://oaidalleapiprodscus.blob.core.windows.net/private/org-lx7MeI8vgrXk2dZIpjAehwph/user-p1Gk7zjjfNg1aoVI1pCcZhUQ/img-JZ6w06yHyjOWuNCLbhxBCARE.png?st=2023-07-22T14%3A33%3A02Z&se=2023-07-22T16%3A33%3A02Z&sp=r&sv=2021-08-06&sr=b&rscd=inline&rsct=image/png&skoid=6aaadede-4fb3-4698-a8f6-684d7786b067&sktid=a48cca56-e6da-484e-a814-9c849652bcb3&skt=2023-07-21T23%3A28%3A09Z&ske=2023-07-22T23%3A28%3A09Z&sks=b&skv=2021-08-06&sig=0bmG0%2BkdTUKFAs2OWVmpBnZLog%2BGy33OPWgDAMbw30Q%3D"})
	//return r, nil
	//
	respo := new(GenImageResp)
	url := "https://api.openai.com/v1/images/generations"
	method := "POST"
	requ := &genImageReq{
		Prompt: prompt,
		N:      2,
		Size:   "1024x1024",
	}
	p, _ := json.Marshal(requ)
	payload := strings.NewReader(string(p))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return respo, "", err
	}
	token := os.Getenv("CHATGPT_TOKEN")

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		return respo, "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return respo, "", err
	}
	mylogrus.MyLog.Infof("genImageBody:%v", string(body))
	err = json.Unmarshal(body, respo)
	return respo, string(body), err
}

func persistentImages(images []string, id mysql.ID) {
	model := domain.CreateModelNil()
	var locals []string
	for _, img := range images {
		dl, err := downloadImage(img)
		if err != nil {
			model.Log.Errorf("downloadImage fail:%v", err)
			continue
		}
		locals = append(locals, dl)
	}
	var image1, image2 string
	for i := range locals {
		if i == 0 {
			image1 = locals[i]
		}
		if i == 1 {
			image2 = locals[i]
		}
	}
	if err := ai_m.UpdateAiImage12(model, id, image1, image2); err != nil {
		model.Log.Errorf("UpdateAiImage12 fail:%v", err)
	}
}

const (
	UPLOAD_PATH = "uploads/file"
	FILE_PORT   = 7002
	DOMAIN      = "http://47.244.34.27"
)

func downloadImage(url string) (string, error) {
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
	}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// 读取文件后缀
	ext := ".png"
	// 读取文件名并md5加密
	// 拼接新文件名(用时间戳)
	filename := time.Now().Format("20060102150405") + ext
	// 创建路径
	err = os.MkdirAll(UPLOAD_PATH, os.ModePerm)
	if err != nil {
		return "", err
	}
	// 拼接路径和文件名
	filepath := UPLOAD_PATH + "/" + filename

	// 创建文件
	out, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// 拷贝文件
	_, err = io.Copy(out, res.Body)
	if err != nil {
		return "", err
	}
	prefix := DOMAIN + fmt.Sprintf(":%d/", FILE_PORT)

	return prefix + filepath, nil

}
