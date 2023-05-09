package chatgpt_r

import (
	"encoding/json"
	"frozen-go-cms/domain/model/chatgpt_m"
	"frozen-go-cms/req"
	"frozen-go-cms/resp"
	"git.hilo.cn/hilo-common/domain"
	"git.hilo.cn/hilo-common/mycontext"
	"git.hilo.cn/hilo-common/resource/mysql"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
)

type ProcessReq struct {
	SessionId uint64           `json:"session_id"`
	Message   []ProcessContent `json:"messages"`
}

type ProcessContent struct {
	Role    string `json:"role"` // user | assistant
	Content string `json:"content"`
}

// @Tags Chatgpt
// @Summary 请求
// @Param Authorization header string true "token"
// @Param ProcessReq body ProcessReq true "请求体"
// @Success 200 {object} ProcessReq
// @Router /v1_0/chatgpt/process [post]
func Process(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	userId, err := req.GetUserId(c)
	if err != nil {
		return myCtx, err
	}
	var param ProcessReq
	if err := c.ShouldBind(&param); err != nil {
		return myCtx, err
	}
	reply, err := process(param)
	if err != nil {
		return myCtx, err
	}
	param.Message = append(param.Message, ProcessContent{Role: "assistant", Content: reply})
	message, _ := json.Marshal(param)
	var model = domain.CreateModelContext(myCtx)
	if err := chatgpt_m.UpdateSessionInit(model, chatgpt_m.ChatgptSession{
		Entity:    mysql.Entity{},
		UserId:    userId,
		SessionId: param.SessionId,
		Message:   string(message),
	}); err != nil {
		return myCtx, err
	}
	resp.ResponseOk(c, param)
	return myCtx, nil
}

func process(param ProcessReq) (string, error) {
	url := "https://jojonull.com/api/chat/process"
	method := "POST"
	c, _ := json.Marshal(param)
	payload := strings.NewReader(string(c))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cookie", "Hm_lvt_71fe5e46112b126abb97279f07e6ac7e=1683509180; connect.sid=s%3AkwnLWmH4HHlY5ubUdMUSv_VREIrh5Mgw.5RxObT6pPbFax0GXNTOuS%2FQkr1eNYhtVFfBIjxGf728; Hm_lpvt_71fe5e46112b126abb97279f07e6ac7e=1683509345")
	req.Header.Add("Origin", "https://jojonull.com")
	req.Header.Add("Referer", "https://jojonull.com/")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36")
	req.Header.Add("sec-ch-ua", "\"Chromium\";v=\"112\", \"Google Chrome\";v=\"112\", \"Not:A-Brand\";v=\"99\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"Windows\"")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), err
}
