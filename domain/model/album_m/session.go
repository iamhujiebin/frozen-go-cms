package album_m

import (
	"git.hilo.cn/hilo-common/domain"
	"git.hilo.cn/hilo-common/resource/mysql"
	"gorm.io/gorm/clause"
)

type Album struct {
	mysql.Entity
	UserId  mysql.ID
	AlbumId mysql.ID
	Content string
}

// 获取用户相册列表
// 一个都没有则初始化一个
func GetUserAlbumsInit(model *domain.Model, userId mysql.ID) ([]Album, error) {
	var Albums []Album
	if err := model.DB().Model(Album{}).Where("user_id = ?", userId).Find(&Albums).Error; err != nil {
		return Albums, err
	}
	if len(Albums) <= 0 {
		if err := model.DB().Model(Album{}).Create(&Album{
			UserId:  userId,
			AlbumId: 0,
		}).Error; err != nil {
			return Albums, err
		}
		if err := model.DB().Model(Album{}).Where("user_id = ?", userId).Find(&Albums).Error; err != nil {
			return Albums, err
		}
	}
	return Albums, nil
}

// 获取用户指定相册
func GetUserAlbum(model *domain.Model, userId, AlbumId mysql.ID) (Album, error) {
	var album Album
	if err := model.DB().Model(Album{}).
		Where("user_id = ? AND album_id = ?", userId, AlbumId).First(&album).Error; err != nil {
		return album, err
	}
	return album, nil
}

// 更新用户相册
func UpdateAlbumInit(model *domain.Model, album Album) error {
	return model.DB().Model(Album{}).
		Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "user_id"}, {Name: "album_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"content": album.Content,
			}),
		}).Create(&album).Error
}

// 更新用户相册
func CreateAlbumInit(model *domain.Model, userId mysql.ID) (mysql.ID, error) {
	var maxAlbum Album
	if err := model.DB().Model(Album{}).Where("user_id = ?", userId).Order("album_id DESC").First(&maxAlbum).Error; err != nil {
		if err := model.DB().Model(Album{}).Create(&maxAlbum).Error; err != nil {
			return 0, err
		}
		return maxAlbum.AlbumId, nil
	}
	maxAlbum = Album{UserId: maxAlbum.UserId, AlbumId: maxAlbum.AlbumId + 1}
	if err := model.DB().Model(Album{}).Create(&maxAlbum).Error; err != nil {
		return 0, err
	}
	return maxAlbum.AlbumId, nil
}

// 删除一个相册
func DeleteAlbum(model *domain.Model, userId, AlbumId mysql.ID) error {
	if err := model.DB().Model(Album{}).Where("user_id = ? AND album_id = ?", userId, AlbumId).Delete(&Album{}).Error; err != nil {
		return err
	}
	return nil
}
