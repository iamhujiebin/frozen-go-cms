package music_r

import (
	"encoding/json"
	"errors"
	"fmt"
	"frozen-go-cms/_const/enum/music_e"
	"frozen-go-cms/common/domain"
	"frozen-go-cms/common/mycontext"
	"frozen-go-cms/common/resource/mysql"
	"frozen-go-cms/domain/model/music_m"
	"frozen-go-cms/resp"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
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

type MusicListReq struct {
	PlaylistId uint64 `form:"playlistId"`
}

// @Tags 音乐
// @Summary 列表
// @Param Authorization header string true "token"
// @Param playlistId query int false "歌单id"
// @Success 200 {object} []Music
// @Router /v1_0/music/list [get]
func MusicList(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	var param MusicListReq
	if err := c.ShouldBindQuery(&param); err != nil {
		return myCtx, err
	}
	var response []Music
	if param.PlaylistId <= 0 {
		musics := music_m.GetAllMusic(model)
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
	} else {
		dbSongs := music_m.GetMusicPlayListSongs(model, param.PlaylistId)
		for _, song := range dbSongs {
			response = append(response, Music{
				Name:   song.Name,
				Artist: song.Artist,
				Url:    song.Url,
				Cover:  song.Cover,
				Lrc:    song.Lrc,
			})
		}
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
	result := search(music_e.SearchTypeSong, param.Q, 100)
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
		Playlists []struct {
			ID          uint64 `json:"id"`
			Name        string `json:"name"`
			CoverImgURL string `json:"coverImgUrl"`
			Creator     struct {
				Nickname  string `json:"nickname"`
				AvatarUrl string `json:"avatarUrl"`
			} `json:"creator"`
			Description string `json:"description"`
			Track       struct {
				Name    string `json:"name"`
				ID      int    `json:"id"`
				Artists []struct {
					Name string `json:"name"`
					ID   int    `json:"id"`
				} `json:"artists"`
				Album struct {
					Name string `json:"name"`
					ID   int    `json:"id"`
				} `json:"album"`
			} `json:"track"`
		} `json:"playlists"`
		HasMore   bool `json:"hasMore"`
		SongCount int  `json:"songCount"`
	} `json:"result"`
	Code int `json:"code"`
}

func search(_type music_e.SearchType, query string, limit int) (response SearchResponse) {
	query = url.QueryEscape(query)
	_url := fmt.Sprintf("http://music.163.com/api/search/get?s=%s&type=%d&limit=%d", query, _type, limit)
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
	songs := down(param.Id, music_e.SongPlayListTypeSong)
	if len(songs) <= 0 {
		return myCtx, errors.New("no song")
	}
	song := songs[0]
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

// @Tags 音乐
// @Summary 删除
// @Param Authorization header string true "token"
// @Param id path int true "歌id"
// @Success 200 {object} Music
// @Router /v1_0/music/{id} [delete]
func MusicDel(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := music_m.DelMusic(model, id); err != nil {
		return myCtx, err
	}
	resp.ResponseOk(c, "")
	return myCtx, nil
}

type DownResponse struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	URL    string `json:"url"`
	Pic    string `json:"pic"`
	Lrc    string `json:"lrc"`
}

func down(id uint64, _type music_e.SongPlayListType) (response []DownResponse) {
	url := fmt.Sprintf("https://api.i-meto.com/meting/api?server=netease&id=%d&type=%s", id, _type)
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
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println(err)
		return
	}
	return response
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

type Playlist struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
	Pic  string `json:"pic"`
}

// @Tags 音乐
// @Summary 搜索歌手歌单
// @Param Authorization header string true "token"
// @Param q  query string true "搜索内容"
// @Success 200 {object} []Playlist
// @Router /v1_0/music/author/search [get]
func MusicAuthorSearch(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	var param MusicSearchReq
	if err := c.ShouldBindQuery(&param); err != nil {
		return myCtx, err
	}
	playlist := search(music_e.SearchTypeAuthor, param.Q, 10)
	var response []Playlist
	for _, list := range playlist.Result.Playlists {
		response = append(response, Playlist{
			Id:   list.ID,
			Name: list.Name,
			Desc: list.Description,
			Pic:  list.Creator.AvatarUrl,
		})
	}
	resp.ResponseOk(c, response)
	return myCtx, nil
}

// @Tags 音乐
// @Summary 歌单列表
// @Param Authorization header string true "token"
// @Success 200 {object} []Playlist
// @Router /v1_0/music/author/list [get]
func MusicAuthorList(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	playlist := music_m.GetAllPlayLists(model)
	var response = []Playlist{{
		Id:   0,
		Name: "默认",
		Desc: "默认歌单",
		Pic:  "",
	}}
	for _, list := range playlist {
		response = append(response, Playlist{
			Id:   list.ID,
			Name: list.Name,
			Desc: list.Desc,
			Pic:  list.Pic,
		})
	}
	resp.ResponseOk(c, response)
	return myCtx, nil
}

type MusicAuthorDownReq struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
	Pic  string `json:"pic"`
}

// @Tags 音乐
// @Summary 下载歌手歌单
// @Param Authorization header string true "token"
// @Param MusicAuthorDownReq body MusicAuthorDownReq true "歌单id"
// @Success 200 {object} []Music
// @Router /v1_0/music/author/down [post]
func MusicAuthorDown(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	var param MusicAuthorDownReq
	if err := c.ShouldBindJSON(&param); err != nil {
		return myCtx, err
	}
	var response []Music
	dbSongs := music_m.GetMusicPlayListSongs(model, param.Id)
	if len(dbSongs) <= 0 {
		songs := down(param.Id, music_e.SongPlayListTypePlaylist)
		for _, song := range songs {
			response = append(response, Music{
				Name:   song.Title,
				Artist: song.Author,
				Url:    song.URL,
				Cover:  song.Pic,
				Lrc:    song.Lrc,
			})
		}
		var playListSongs []music_m.MusicPlaylistSongs
		for i, song := range response {
			response[i].Lrc = downloadLrc(song.Lrc)
			playListSongs = append(playListSongs, music_m.MusicPlaylistSongs{
				PlaylistId: param.Id,
				Name:       song.Name,
				Artist:     song.Artist,
				Url:        song.Url,
				Cover:      song.Cover,
				Lrc:        response[i].Lrc,
			})
		}
		if err := music_m.AddMusicPlayListSongs(model, param.Id, param.Name, param.Desc, param.Pic, playListSongs); err != nil {
			return myCtx, err
		}
	} else {
		for _, song := range dbSongs {
			response = append(response, Music{
				Name:   song.Name,
				Artist: song.Artist,
				Url:    song.Url,
				Cover:  song.Cover,
				Lrc:    song.Lrc,
			})
		}
	}
	resp.ResponseOk(c, response)
	return myCtx, nil
}
