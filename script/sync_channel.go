package main

import (
	"encoding/json"
	"fmt"
	"frozen-go-cms/domain/model/channel_m"
	"git.hilo.cn/hilo-common/domain"
	"git.hilo.cn/hilo-common/resource/mysql"
	"io/ioutil"
	"net/http"
)

func main() {

	url := "http://geek.itheima.net/v1_0/channels"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJodHRwczovL3d3dy5pdGNhc3QuY24vIiwic3ViIjoiMTExMSIsImp0aSI6IjM4NmQwZWFlLTMzNDEtNDliOC05OWE4LWE4NTBkNTM4Y2MyZCIsImlhdCI6MTY4MzE1Mjk2NSwiZXhwIjoxNjgzMTU2NTY1fQ.o5vej5gPtJ2KytvC3v2OXh4v4fp_Be2oardmew0RKmI")

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
	type Response struct {
		Data struct {
			Channels []struct {
				Id   uint64 `json:"id"`
				Name string `json:"name"`
			} `json:"channels"`
		} `json:"data"`
	}
	data := new(Response)
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}
	model := domain.CreateModelNil()
	for _, v := range data.Data.Channels {
		if v.Id == 0 {
			continue
		}
		if err := channel_m.AddChannel(model, channel_m.Channel{
			Entity: mysql.Entity{
				ID: v.Id,
			},
			Name: v.Name,
		}); err != nil {
			panic(err)
		}
	}
}
