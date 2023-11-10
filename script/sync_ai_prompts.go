package main

import (
	"encoding/json"
	"fmt"
	"frozen-go-cms/common/domain"
	"frozen-go-cms/domain/model/ai_m"
	"io/ioutil"
	"net/http"
)

type AiPromptResp []AiPromptData

type AiPromptData struct {
	Code string `json:"code"`
	Name struct {
		En string `json:"en"`
		Zh string `json:"zh"`
	} `json:"name"`
	SubTabs []struct {
		Code string `json:"code"`
		Name struct {
			En string `json:"en"`
			Zh string `json:"zh"`
		} `json:"name"`
		Prompts []struct {
			En   string `json:"en"`
			Zh   string `json:"zh"`
			Code string `json:"code"`
		} `json:"prompts"`
	} `json:"subTabs"`
}

func main() {

	url := "https://image.gptplus.io/userspace/settings/remaker.ai.web.tabsMarket.v1.json"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("sec-ch-ua", "\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", \"Google Chrome\";v=\"114\"")
	req.Header.Add("Referer", "https://remaker.ai/")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	req.Header.Add("sec-ch-ua-platform", "\"Windows\"")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var resp []AiPromptData
	json.Unmarshal(body, &resp)
	model := domain.CreateModelNil()
	idByCode := ai_m.GetIdByCode(model)
	for _, v := range resp {
		for _, v2 := range v.SubTabs {
			parentId := idByCode[v2.Code]
			for _, v3 := range v2.Prompts {
				if err := ai_m.AddAiPrompt(model, ai_m.AiPrompt{
					Zh:       v3.Zh,
					En:       v3.En,
					Level:    2,
					ParentId: uint64(parentId),
					Code:     v3.Code,
				}); err != nil {
					panic(err)
				}
			}
		}
	}
}
