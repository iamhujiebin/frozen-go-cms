package discord_m

import (
	"frozen-go-cms/common/domain"
	"frozen-go-cms/common/resource/mysql"
)

// DiscordFeed  discord抓包
type DiscordFeed struct {
	mysql.Entity
	MsgId       string `gorm:"column:msg_id" json:"msg_id"`             //  msg_id
	MsgContent  string `gorm:"column:msg_content" json:"msg_content"`   //  消息内容
	Author      string `gorm:"column:author" json:"author"`             //  作者名
	Avatar      string `gorm:"column:avatar" json:"avatar"`             //  作者头像
	AttachId    string `gorm:"column:attach_id" json:"attach_id"`       //  attach_id
	Url         string `gorm:"column:url" json:"url"`                   //  图片1
	ProxyUrl    string `gorm:"column:proxy_url" json:"proxy_url"`       //  图片2
	Width       int    `gorm:"column:width" json:"width"`               //  width
	Height      int    `gorm:"column:height" json:"height"`             //  height
	Size        int    `gorm:"column:size" json:"size"`                 //  size
	ContentType string `gorm:"column:content_type" json:"content_type"` //  content_type
}

func (DiscordFeed) TableName() string {
	return "discord_feed"
}

func AddDiscordFeed(model *domain.Model, feed DiscordFeed) error {
	return model.DB().Create(&feed).Error
}
