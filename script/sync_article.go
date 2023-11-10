package main

import (
	"encoding/json"
	"fmt"
	"frozen-go-cms/common/domain"
	"frozen-go-cms/common/resource/mysql"
	"frozen-go-cms/domain/model/article_m"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {

	url := "http://geek.itheima.net/v1_0/mp/articles?page=1&per_page=10000&channel_id="
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
			Result []struct {
				Id           uint64    `json:"id,string"`
				Title        string    `json:"title"`
				Status       int       `json:"status"`
				CommentCount mysql.Num `json:"comment_count"`
				PubDate      string    `json:"pubdate"`
				Cover        struct {
					Type   int      `json:"type"`
					Images []string `json:"images"`
				} `json:"cover"`
				LikeCount mysql.Num `json:"like_count"`
				ReadCount mysql.Num `json:"read_count"`
			} `json:"results"`
		} `json:"data"`
	}
	data := new(Response)
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}
	model := domain.CreateModelNil()
	for _, v := range data.Data.Result {
		url := "http://geek.itheima.net/v1_0/mp/articles/" + fmt.Sprintf("%d", v.Id)
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
				ChannelId mysql.ID `json:"channel_id"`
				Content   string   `json:"content"`
			} `json:"data"`
		}
		rr := new(Response)
		err = json.Unmarshal(body, &rr)
		if err != nil {
			panic(err)
		}
		images, err := json.Marshal(v.Cover.Images)
		t, _ := time.Parse("2006-01-02 15:04:05", v.PubDate)
		if _, err := article_m.AddArticle(model, article_m.Article{
			Entity: mysql.Entity{
				ID: v.Id,
			},
			ChannelId:    rr.Data.ChannelId,
			Title:        v.Title,
			Status:       v.Status,
			CommentCount: v.CommentCount,
			LikeCount:    v.LikeCount,
			ReadCount:    v.ReadCount,
			PubDate:      t,
			CoverType:    v.Cover.Type,
			CoverImages:  string(images),
			Content:      rr.Data.Content,
		}); err != nil {
			panic(err)
		}
	}
}
