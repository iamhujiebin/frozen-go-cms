package music_m

import (
	"git.hilo.cn/hilo-common/domain"
	"git.hilo.cn/hilo-common/resource/mysql"
)

// Music  music列表
type Music struct {
	mysql.Entity
	Name   string `gorm:"column:name"`   //  歌曲名
	Artist string `gorm:"column:artist"` //  歌手
	Url    string `gorm:"column:url"`    //  音乐mp3
	Cover  string `gorm:"column:cover"`  //  音乐封面
	Lrc    string `gorm:"column:lrc"`    //  歌词
}

func (Music) TableName() string {
	return "music"
}

func AddMusic(model *domain.Model, music Music) error {
	return model.DB().Create(&music).Error
}

func GetAllMusic(model *domain.Model) []Music {
	var musics []Music
	if err := model.DB().Model(Music{}).Find(&musics).Error; err != nil {
		model.Log.Errorf("GetAllMusic fail:%v", err)
	}
	return musics
}
