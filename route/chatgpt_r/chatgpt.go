package chatgpt_r

import (
	"encoding/json"
	"errors"
	"frozen-go-cms/_const/enum/ws_e"
	"frozen-go-cms/domain/model/chatgpt_m"
	"frozen-go-cms/req"
	"frozen-go-cms/resp"
	"frozen-go-cms/route/ws_r"
	"git.hilo.cn/hilo-common/domain"
	"git.hilo.cn/hilo-common/mycontext"
	"git.hilo.cn/hilo-common/resource/mysql"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type ProcessReq struct {
	SessionId uint64           `json:"session_id"`
	Message   []ProcessContent `json:"messages"`
}

type ProcessContent struct {
	Role        string `json:"role"` // user | assistant
	Content     string `json:"content"`
	CreatedTime string `json:"createdTime"`
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
	reply, err := RealProcess(param)
	if err != nil {
		return myCtx, err
	}
	if len(param.Message) > 0 {
		param.Message[len(param.Message)-1].CreatedTime = time.Now().Format("2006-01-02 15:04:05")
	}
	param.Message = append(param.Message, ProcessContent{Role: "assistant", Content: reply, CreatedTime: time.Now().Format("2006-01-02 15:04:05")})
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
	ws_r.SendToClient(userId, ws_e.CmdNewMsg) // 多端同步
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

type ChatGPTRequest struct {
	Model       string    `json:"model"` // gpt-3.5-turbo
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"` // 0.7
	N           int       `json:"n"`           // 返回答案个数
	Stream      bool      `json:"stream"`      // 是否流式返回
}

type ChatGPTResponse struct {
	ID      string    `json:"id"`
	Object  string    `json:"object"`
	Created int       `json:"created"`
	Model   string    `json:"model"`
	Usage   Usage     `json:"usage"`
	Choices []Choices `json:"choices"`
}
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type Choices struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
}

func RealProcess(p ProcessReq) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"
	method := "POST"

	var param = ChatGPTRequest{
		Model:       "gpt-3.5-turbo",
		Messages:    nil,
		Temperature: 0.7,
		N:           1,
		Stream:      false,
	}
	for _, message := range p.Message {
		if message.Role == "assistant" {
			continue
		}
		param.Messages = append(param.Messages, Message{message.Role, message.Content})
	}
	j, _ := json.Marshal(param)
	payload := strings.NewReader(string(j))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return "", err
	}
	token := os.Getenv("CHATGPT_TOKEN")

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	response := new(ChatGPTResponse)
	err = json.Unmarshal(body, response)
	if err != nil {
		return "", err
	}
	if len(response.Choices) <= 0 {
		return "", errors.New("chatgpt no answers")
	}
	return response.Choices[0].Message.Content, nil
}
