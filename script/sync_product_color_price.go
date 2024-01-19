package main

import (
	"encoding/json"
	"fmt"
	"frozen-go-cms/common/domain"
	"frozen-go-cms/domain/model/product_price_m"
	"github.com/spf13/cast"
	"io/ioutil"
	"net/http"
)

type AutoGenerate3 struct {
	Status   string `json:"status,omitempty"`
	Page     string `json:"page,omitempty"`
	Records  string `json:"records,omitempty"`
	Total    string `json:"total,omitempty"`
	KWord    string `json:"KWord,omitempty"`
	Message  string `json:"message,omitempty"`
	DataList []struct {
		Rownum               string `json:"rownum,omitempty"`
		Pid                  string `json:"PID,omitempty"`
		EColorlName          string `json:"E_ColorlName,omitempty"`
		ECreateUser          string `json:"E_CreateUser,omitempty"`
		ECreateDate          string `json:"E_CreateDate,omitempty"`
		EUpdateUser          string `json:"E_UpdateUser,omitempty"`
		EUpdateDate          string `json:"E_UpdateDate,omitempty"`
		ECreateIP            string `json:"E_CreateIP,omitempty"`
		EState               string `json:"E_State,omitempty"`
		EIndex               string `json:"E_Index,omitempty"`
		EPguid               string `json:"E_pguid,omitempty"`
		EColorCode           string `json:"E_ColorCode,omitempty"`
		EPrintStartNum       string `json:"E_PrintStartNum,omitempty"`
		EPrintStartPrice     string `json:"E_PrintStartPrice,omitempty"`
		EPrintBaseNum        string `json:"E_PrintBaseNum,omitempty"`
		EPrintBasePrice      string `json:"E_PrintBasePrice,omitempty"`
		EPrintLowSumNum      string `json:"E_PrintLowSumNum,omitempty"`
		EPrintLowNum         string `json:"E_PrintLowNum,omitempty"`
		EPrintBetweenNum     string `json:"E_PrintBetweenNum,omitempty"`
		EPrintBetweenAddNum  string `json:"E_PrintBetweenAddNum,omitempty"`
		EISCover             string `json:"E_ISCover,omitempty"`
		EISPage1             string `json:"E_ISPage1,omitempty"`
		EISPage2             string `json:"E_ISPage2,omitempty"`
		EISCard              string `json:"E_ISCard,omitempty"`
		EISBox               string `json:"E_ISBox,omitempty"`
		EPrintStartNum2      string `json:"E_PrintStartNum2,omitempty"`
		EPrintStartPrice2    string `json:"E_PrintStartPrice2,omitempty"`
		EPrintBaseNum2       string `json:"E_PrintBaseNum2,omitempty"`
		EPrintBasePrice2     string `json:"E_PrintBasePrice2,omitempty"`
		EPrintLowSumNum2     string `json:"E_PrintLowSumNum2,omitempty"`
		EPrintLowNum2        string `json:"E_PrintLowNum2,omitempty"`
		EPrintBetweenNum2    string `json:"E_PrintBetweenNum2,omitempty"`
		EPrintBetweenAddNum2 string `json:"E_PrintBetweenAddNum2,omitempty"`
	} `json:"DataList,omitempty"`
}

func main() {

	url := "http://zf.jenyun.com/OfferPrintColorList/DataListGet?KWord=&_search=false&nd=1705681567692&rows=500&page=1&sidx=&sord=asc"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Cookie", "ASP.NET_SessionId=pwi0uasq0xrensyto4qulgmd; OfferMaterialListIndex=1; OfferCraftListIndex=1; OfferPrintColorListIndex=1; ConrtactConfigIndex=1; OfferSizeListIndex=1; SysConfigIndex=1; MobileMoneyIndex=1; MobileCustomerIndex=1; MobileSuppliersIndex=1; MaterialListIndex=1")
	req.Header.Add("Referer", "http://zf.jenyun.com/OfferPrintColorList")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

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
	var response AutoGenerate3
	json.Unmarshal(body, &response)
	model := domain.CreateModelNil()
	for _, v := range response.DataList {
		err := product_price_m.CreateColorPrice(model, product_price_m.ColorPrice{
			ColorName:           v.EColorlName,
			ColorCode:           v.EColorCode,
			PrintStartNum:       cast.ToInt64(v.EPrintStartNum),
			PrintStartPrice:     cast.ToFloat64(v.EPrintStartPrice),
			PrintBaseNum:        cast.ToInt64(v.EPrintBaseNum),
			PrintBasePrice:      cast.ToFloat64(v.EPrintBasePrice),
			PrintLowSumNum:      cast.ToInt64(v.EPrintLowSumNum),
			PrintLowNum:         cast.ToInt64(v.EPrintLowNum),
			PrintBetweenNum:     cast.ToInt64(v.EPrintBetweenNum),
			PrintBetweenAddNum:  cast.ToInt64(v.EPrintBetweenAddNum),
			PrintStartNum2:      cast.ToInt64(v.EPrintStartNum2),
			PrintStartPrice2:    cast.ToFloat64(v.EPrintStartPrice2),
			PrintBaseNum2:       cast.ToInt64(v.EPrintBaseNum2),
			PrintBasePrice2:     cast.ToFloat64(v.EPrintBasePrice2),
			PrintLowSumNum2:     cast.ToInt64(v.EPrintLowSumNum2),
			PrintLowNum2:        cast.ToInt64(v.EPrintLowNum2),
			PrintBetweenNum2:    cast.ToInt64(v.EPrintBetweenNum2),
			PrintBetweenAddNum2: cast.ToInt64(v.EPrintBetweenAddNum2),
			PageCover:           cast.ToInt64(v.EISCover),
			PageInner:           cast.ToInt64(v.EISPage1),
			PageTag:             cast.ToInt64(v.EISPage2),
			Card:                cast.ToInt64(v.EISCard),
			Box:                 cast.ToInt64(v.EISBox),
			Index:               cast.ToInt64(v.EIndex),
			CreateIp:            v.ECreateIP,
			CreateUser:          v.ECreateUser,
			UpdateUser:          v.EUpdateUser,
		})
		if err != nil {
			panic(err)
		}
	}
	fmt.Println(string(body))
}
