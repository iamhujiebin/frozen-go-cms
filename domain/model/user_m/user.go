package user_m

import (
	"git.hilo.cn/hilo-common/domain"
	"git.hilo.cn/hilo-common/resource/mysql"
	"gorm.io/gorm"
)

type User struct {
	mysql.Entity
	Mobile string
}

// 获取用户
func GetUser(model *domain.Model, id uint64) User {
	var user User
	if err := model.DB().Model(User{}).Where("id = ?", id).First(&user).Error; err != nil {
		model.Log.Errorf("GetUser err:%v", err)
	}
	return user
}

// 获取或者创建用户
func GetUserOrCreate(model *domain.Model, mobile string) (User, error) {
	var user User
	if err := model.DB().Model(User{}).Where("mobile = ?", mobile).First(&user).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return User{}, err
		}
		user.Mobile = mobile
		if err := model.DB().Model(User{}).Create(&user).Error; err != nil {
			return User{}, err
		}
	}
	return user, nil
}
