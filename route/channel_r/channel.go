package channel_r

import (
	"frozen-go-cms/domain/model/channel_m"
	"frozen-go-cms/resp"
	"git.hilo.cn/hilo-common/domain"
	"git.hilo.cn/hilo-common/mycontext"
	"github.com/gin-gonic/gin"
)

type ChannelData struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

type ChannelsResp struct {
	Channels []ChannelData `json:"channels"`
}

// @Tags 频道
// @Summary 列表
// @Success 200 {object} ChannelsResp
// @Router /v1_0/channels [get]
func Channels(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	channels := channel_m.GetChannels(model)
	var response ChannelsResp
	for _, channel := range channels {
		response.Channels = append(response.Channels, ChannelData{
			Id:   channel.ID,
			Name: channel.Name,
		})
	}
	resp.ResponseOk(c, response)
	return myCtx, nil
}
