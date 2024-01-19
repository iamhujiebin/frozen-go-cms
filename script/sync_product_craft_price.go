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

type AutoGenerate2 struct {
	Status   string `json:"status,omitempty"`
	Page     string `json:"page,omitempty"`
	Records  string `json:"records,omitempty"`
	Total    string `json:"total,omitempty"`
	KWord    string `json:"KWord,omitempty"`
	Message  string `json:"message,omitempty"`
	DataList []struct {
		Rownum            string `json:"rownum,omitempty"`
		Pid               string `json:"PID,omitempty"`
		ECraftBobyName    string `json:"E_CraftBobyName,omitempty"`
		ECraftName        string `json:"E_CraftName,omitempty"`
		ECraftPrice       string `json:"E_CraftPrice,omitempty"`
		EImg1             string `json:"E_Img1,omitempty"`
		ECreateUser       string `json:"E_CreateUser,omitempty"`
		ECreateDate       string `json:"E_CreateDate,omitempty"`
		EUpdateUser       string `json:"E_UpdateUser,omitempty"`
		EUpdateDate       string `json:"E_UpdateDate,omitempty"`
		ECreateIP         string `json:"E_CreateIP,omitempty"`
		EState            string `json:"E_State,omitempty"`
		EIndex            string `json:"E_Index,omitempty"`
		EPguid            string `json:"E_pguid,omitempty"`
		ECraftBobyCode    string `json:"E_CraftBobyCode,omitempty"`
		ECraftCode        string `json:"E_CraftCode,omitempty"`
		ECraftMinSumPrice string `json:"E_CraftMinSumPrice,omitempty"`
		ECraftModelPrice  string `json:"E_CraftModelPrice,omitempty"`
		ECraftUnit        string `json:"E_CraftUnit,omitempty"`
		EBookPrice        string `json:"E_BookPrice,omitempty"`
		ETaskName         string `json:"E_TaskName,omitempty"`
	} `json:"DataList,omitempty"`
}

func main() {

	url := "http://zf.jenyun.com/OfferCraftList/DataListGet?KWord=&_search=false&nd=1705681473527&rows=500&page=1&sidx=&sord=asc"
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
	req.Header.Add("Referer", "http://zf.jenyun.com/OfferCraftList")
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
	var response AutoGenerate2
	json.Unmarshal(body, &response)
	model := domain.CreateModelNil()
	for _, v := range response.DataList {
		err := product_price_m.CreateCraftPrice(model, product_price_m.CraftPrice{
			CraftName:       v.ECraftName,
			CraftCode:       v.ECraftCode,
			CraftBodyName:   v.ECraftBobyName,
			CraftBodyCode:   v.ECraftBobyCode,
			CraftPrice:      cast.ToFloat64(v.ECraftPrice),
			BookPrice:       cast.ToFloat64(v.EBookPrice),
			CraftUnit:       v.ECraftUnit,
			MinSumPrice:     cast.ToFloat64(v.ECraftMinSumPrice),
			CraftModelPrice: cast.ToFloat64(v.ECraftModelPrice),
			TaskName:        v.ETaskName,
			Index:           cast.ToInt64(v.EIndex),
			CreateIp:        v.ECreateIP,
			CreateUser:      v.ECreateUser,
			UpdateUser:      v.EUpdateUser,
		})
		if err != nil {
			panic(err)
		}
	}
	fmt.Println(string(body))
}
