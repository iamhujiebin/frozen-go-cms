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

type AutoGenerate4 struct {
	Status   string `json:"status,omitempty"`
	Page     string `json:"page,omitempty"`
	Records  string `json:"records,omitempty"`
	Total    string `json:"total,omitempty"`
	KWord    string `json:"KWord,omitempty"`
	Message  string `json:"message,omitempty"`
	DataList []struct {
		Rownum             string `json:"rownum,omitempty"`
		Pid                string `json:"PID,omitempty"`
		ESizeName          string `json:"E_SizeName,omitempty"`
		ESizeWidth         string `json:"E_SizeWidth,omitempty"`
		ESizeHeight        string `json:"E_SizeHeight,omitempty"`
		ESizePrice         string `json:"E_SizePrice,omitempty"`
		ESizeOpenNum       string `json:"E_SizeOpenNum,omitempty"`
		ECreateUser        string `json:"E_CreateUser,omitempty"`
		ECreateDate        string `json:"E_CreateDate,omitempty"`
		EUpdateUser        string `json:"E_UpdateUser,omitempty"`
		EUpdateDate        string `json:"E_UpdateDate,omitempty"`
		ECreateIP          string `json:"E_CreateIP,omitempty"`
		EState             string `json:"E_State,omitempty"`
		EIndex             string `json:"E_Index,omitempty"`
		EPguid             string `json:"E_pguid,omitempty"`
		ESizeCode          string `json:"E_SizeCode,omitempty"`
		ESizeWidthMax      string `json:"E_SizeWidthMax,omitempty"`
		ESizeWidthMin      string `json:"E_SizeWidthMin,omitempty"`
		ESizeHeightMax     string `json:"E_SizeHeightMax,omitempty"`
		ESizeHeightMin     string `json:"E_SizeHeightMin,omitempty"`
		EPerSqmX           string `json:"E_PerSqmX,omitempty"`
		EPerSqmY           string `json:"E_PerSqmY,omitempty"`
		EDeviceWidth       string `json:"E_DeviceWidth,omitempty"`
		EDeviceHight       string `json:"E_DeviceHight,omitempty"`
		EDeviceAddBase     string `json:"E_DeviceAddBase,omitempty"`
		EDeviceAddPosition string `json:"E_DeviceAddPosition,omitempty"`
		EType              string `json:"E_Type,omitempty"`
	} `json:"DataList,omitempty"`
}

func main() {

	url := "http://zf.jenyun.com/OfferSizeList/DataListGet?KWord=&FType=&_search=false&nd=1705683124208&rows=50&page=1&sidx=&sord=asc"
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
	req.Header.Add("Referer", "http://zf.jenyun.com/OfferSizeList")
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
	var response AutoGenerate4
	json.Unmarshal(body, &response)
	model := domain.CreateModelNil()
	for _, v := range response.DataList {
		err := product_price_m.CreateSizeConfig(model, product_price_m.SizeConfig{
			SizeName:          v.ESizeName,
			SizeCode:          v.ESizeCode,
			Type:              cast.ToInt64(v.EType),
			SizeWidth:         cast.ToInt64(v.ESizeWidth),
			SizeWidthMax:      cast.ToInt64(v.ESizeWidthMax),
			SizeWidthMin:      cast.ToInt64(v.ESizeWidthMin),
			SizeHeight:        cast.ToInt64(v.ESizeHeight),
			SizeHeightMax:     cast.ToInt64(v.ESizeHeightMax),
			SizeHeightMin:     cast.ToInt64(v.ESizeHeightMin),
			PerSqmX:           cast.ToFloat64(v.EPerSqmX),
			PerSqmY:           cast.ToFloat64(v.EPerSqmY),
			DeviceWidth:       cast.ToInt64(v.EDeviceWidth),
			DeviceHeight:      cast.ToInt64(v.EDeviceHight),
			DeviceAddBase:     cast.ToInt64(v.EDeviceAddBase),
			DeviceAddPosition: cast.ToInt64(v.EDeviceAddPosition),
			SizeOpenNum:       cast.ToInt64(v.ESizeOpenNum),
			Index:             cast.ToInt64(v.EIndex),
			CreateIp:          v.ECreateIP,
			CreateUser:        v.ECreateUser,
			UpdateUser:        v.EUpdateUser,
		})
		if err != nil {
			panic(err)
		}
	}
	fmt.Println(string(body))
}
