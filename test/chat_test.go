package test

import (
	"frozen-go-cms/route/chatgpt_r"
	"testing"
)

func TestChatGPT(t *testing.T) {
	res, err := chatgpt_r.RealProcess(chatgpt_r.ProcessReq{
		SessionId: 0,
		Message: []chatgpt_r.ProcessContent{
			{
				Role:    "user",
				Content: "你好",
			},
		},
	})
	if err != nil {
		println(err.Error())
	}
	println(res)
}
