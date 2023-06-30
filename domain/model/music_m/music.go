package music_m

import (
	"frozen-go-cms/hilo-common/domain"
	"frozen-go-cms/hilo-common/resource/mysql"
	"gorm.io/gorm/clause"
	"strings"
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

// MusicPlaylist  歌单列表
type MusicPlaylist struct {
	mysql.Entity
	Name string `gorm:"column:name"` //  歌单名
	Desc string `gorm:"column:desc"` //  描述
	Pic  string `gorm:"column:pic"`  //  图片
}

func (MusicPlaylist) TableName() string {
	return "music_playlist"
}

// MusicPlaylistSongs  歌单歌曲列表
type MusicPlaylistSongs struct {
	mysql.Entity
	PlaylistId uint64 `gorm:"column:playlist_id"` //  歌单id
	Name       string `gorm:"column:name"`        //  歌曲名
	Artist     string `gorm:"column:artist"`      //  歌手
	Url        string `gorm:"column:url"`         //  音乐mp3
	Cover      string `gorm:"column:cover"`       //  音乐封面
	Lrc        string `gorm:"column:lrc"`         //  歌词
}

func (MusicPlaylistSongs) TableName() string {
	return "music_playlist_songs"
}

func AddMusic(model *domain.Model, music Music) error {
	return model.DB().Clauses(clause.OnConflict{UpdateAll: true}).Create(&music).Error
}

func GetAllMusic(model *domain.Model) []Music {
	var musics []Music
	if err := model.DB().Model(Music{}).Order("updated_time DESC").Find(&musics).Error; err != nil {
		model.Log.Errorf("GetAllMusic fail:%v", err)
	}
	return musics
}

func AddMusicPlayListSongs(model *domain.Model, playListId mysql.ID, name, desc, pic string, songs []MusicPlaylistSongs) error {
	return model.Transaction(func(model *domain.Model) error {
		if err := model.DB().Create(&MusicPlaylist{
			Entity: mysql.Entity{ID: playListId},
			Name:   name,
			Desc:   desc,
			Pic:    pic,
		}).Error; err != nil {
			if strings.Contains(err.Error(), "Duplicate") { // 之前下载过的歌单就不下载
				return nil
			}
			return err
		}
		return model.DB().Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(songs, 1000).Error
	})
}

func GetMusicPlayListSongs(model *domain.Model, playListId mysql.ID) []MusicPlaylistSongs {
	var res []MusicPlaylistSongs
	if err := model.DB().Model(MusicPlaylistSongs{}).Where("playlist_id = ?", playListId).Find(&res).Error; err != nil {
		model.Log.Errorf("GetMusicPlayListSongs :%v", err)
		return res
	}
	return res
}

func GetAllPlayLists(model *domain.Model) []MusicPlaylist {
	var res []MusicPlaylist
	if err := model.DB().Model(MusicPlaylist{}).Order("created_time").Find(&res).Error; err != nil {
		model.Log.Errorf("GetAllPlayLists :%v", err)
		return res
	}
	return res
}

func DelMusic(model *domain.Model, id mysql.ID) error {
	return model.DB().Model(Music{}).Where("id = ?", id).Delete(&Music{}).Error
}
