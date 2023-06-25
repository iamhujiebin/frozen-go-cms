package main

import (
	"encoding/json"
	"fmt"
	"frozen-go-cms/domain/model/music_m"
	"frozen-go-cms/hilo-common/domain"
	"io/ioutil"
	"net/http"
)

type Data struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
	Url    string `json:"url"`
	Cover  string `json:"cover"`
	Lrc    string `json:"lrc"`
	LrcStr string `json:"lrcStr"`
}

func main() {

	url := "https://www.noiseblog.top/json/music.json"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
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
	fmt.Println(string(body))
	var datas []Data
	if err := json.Unmarshal(body, &datas); err != nil {
		panic(err)
	}
	for i, v := range datas {
		datas[i].LrcStr = downloadLrc(v.Lrc)
		if err := music_m.AddMusic(domain.CreateModelNil(), music_m.Music{
			Name:   v.Name,
			Artist: v.Artist,
			Url:    v.Url,
			Cover:  v.Cover,
			Lrc:    datas[i].LrcStr,
		}); err != nil {
			panic(err)
		}
	}
}

func downloadLrc(url string) string {

	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
	}
	res, err := client.Do(req)
	if err != nil {
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	return string(body)
}
