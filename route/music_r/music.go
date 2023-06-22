package music_r

import (
	"encoding/json"
	"fmt"
	"frozen-go-cms/domain/model/music_m"
	"frozen-go-cms/resp"
	"git.hilo.cn/hilo-common/domain"
	"git.hilo.cn/hilo-common/mycontext"
	"git.hilo.cn/hilo-common/resource/mysql"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Music struct {
	Id       uint64 `json:"id"`
	Name     string `json:"name"`     //  歌曲名
	Artist   string `json:"artist"`   //  歌手
	Url      string `json:"url"`      //  音乐mp3
	Cover    string `json:"cover"`    //  音乐封面
	Lrc      string `json:"lyric"`    //  歌词
	Duration string `json:"duration"` // 格式化歌曲时长
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
			Id:       music.ID,
			Name:     music.Name,
			Artist:   artist,
			Cover:    music.Album.Artist.Img1V1URL,
			Duration: formatTime(music.Duration),
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
			Duration uint64 `json:"duration"`
		} `json:"songs"`
		HasMore   bool `json:"hasMore"`
		SongCount int  `json:"songCount"`
	} `json:"result"`
	Code int `json:"code"`
}

func search(query string) (response SearchResponse) {
	query = url.QueryEscape(query)
	_url := fmt.Sprintf("http://music.163.com/api/search/get?s=%s&type=1&limit=100", query)
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

type MusicDownReq struct {
	Id uint64 `form:"id"`
}

// @Tags 音乐
// @Summary 搜索
// @Param Authorization header string true "token"
// @Param id query int true "歌id"
// @Success 200 {object} Music
// @Router /v1_0/music/down [get]
func MusicDown(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	var param MusicDownReq
	if err := c.ShouldBindQuery(&param); err != nil {
		return myCtx, err
	}
	song := down(param.Id)
	var response = Music{
		Id:     param.Id,
		Name:   song.Title,
		Artist: song.Author,
		Url:    song.URL,
		Cover:  song.Pic,
		Lrc:    song.Lrc,
	}
	response.Lrc = downloadLrc(response.Lrc)
	if err := music_m.AddMusic(model, music_m.Music{
		Entity: mysql.Entity{ID: response.Id},
		Name:   response.Name,
		Artist: response.Artist,
		Url:    response.Url,
		Cover:  response.Cover,
		Lrc:    response.Lrc,
	}); err != nil {
		return myCtx, err
	}
	resp.ResponseOk(c, response)
	return myCtx, nil
}

type DownResponse struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	URL    string `json:"url"`
	Pic    string `json:"pic"`
	Lrc    string `json:"lrc"`
}

func down(id uint64) (response DownResponse) {
	url := fmt.Sprintf("https://api.i-meto.com/meting/api?server=netease&id=%d&type=song", id)
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
	var responses []DownResponse
	if err := json.Unmarshal(body, &responses); err != nil {
		fmt.Println(err)
		return
	}
	if len(responses) <= 0 {
		return
	}
	return responses[0]
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

func formatTime(milliSeconds uint64) string {
	duration := time.Duration(milliSeconds) * time.Millisecond
	return fmt.Sprintf("%02d:%02d\n", int(duration.Minutes())%60, int(duration.Seconds())%60)
}
