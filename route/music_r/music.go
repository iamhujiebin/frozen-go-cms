package music_r

import (
	"frozen-go-cms/domain/model/music_m"
	"frozen-go-cms/resp"
	"git.hilo.cn/hilo-common/domain"
	"git.hilo.cn/hilo-common/mycontext"
	"github.com/gin-gonic/gin"
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
