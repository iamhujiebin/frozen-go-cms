package music_r

import (
	"encoding/json"
	"fmt"
	"frozen-go-cms/domain/model/music_m"
	"frozen-go-cms/resp"
	"git.hilo.cn/hilo-common/domain"
	"git.hilo.cn/hilo-common/mycontext"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Music struct {
	Id     uint64 `json:"id"`
	Name   string `json:"name"`   //  歌曲名
	Artist string `json:"artist"` //  歌手
	Url    string `json:"url"`    //  音乐mp3
	Cover  string `json:"cover"`  //  音乐封面
	Lrc    string `json:"lrc"`    //  歌词
}

// @Tags 音乐
// @Summary 列表
// @Param Authorization header string true "token"
// @Success 200 {object} []Music
// @Router /v1_0/music/list [get]
func MusicList(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	musics := music_m.GetAllMusic(model)
	var response []Music
	for _, music := range musics {
		response = append(response, Music{
			Id:     music.ID,
			Name:   music.Name,
			Artist: music.Artist,
			Url:    music.Url,
			Cover:  music.Cover,
			Lrc:    music.Lrc,
		})
	}
	resp.ResponseOk(c, response)
	return myCtx, nil
}

type MusicSearchReq struct {
	Q string `form:"q"`
}

// @Tags 音乐
// @Summary 搜索
// @Param Authorization header string true "token"
// @Param q query string true "搜索词"
// @Success 200 {object} []Music
// @Router /v1_0/music/search [get]
func MusicSearch(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	var param MusicSearchReq
	if err := c.ShouldBindQuery(&param); err != nil {
		return myCtx, err
	}
	result := search(param.Q)
	var response []Music
	for _, music := range result.Result.Songs {
		artist := ""
		if len(music.Artists) > 0 {
			artist = music.Artists[0].Name
		}
		response = append(response, Music{
			Id:     music.ID,
			Name:   music.Name,
			Artist: artist,
			Cover:  music.Album.Artist.Img1V1URL,
		})
	}
	resp.ResponseOk(c, response)
	return myCtx, nil
}

type SearchResponse struct {
	Result struct {
		Songs []struct {
			ID      uint64 `json:"id"`
			Name    string `json:"name"`
			Artists []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"artists"`
			Album struct {
				ID     int    `json:"id"`
				Name   string `json:"name"`
				Artist struct {
					ID        int    `json:"id"`
					Name      string `json:"name"`
					Img1V1URL string `json:"img1v1Url"`
				} `json:"artist"`
			} `json:"album"`
		} `json:"songs"`
		HasMore   bool `json:"hasMore"`
		SongCount int  `json:"songCount"`
	} `json:"result"`
	Code int `json:"code"`
}

func search(query string) (response SearchResponse) {
	query = url.QueryEscape(query)
	_url := fmt.Sprintf("http://music.163.com/api/search/get?s=%s&type=1", query)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, _url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Referer", "https://music.163.com")

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
	if err := json.Unmarshal(body, &response); err != nil {
		return
	}
	return
}
