package user_m

import (
	"frozen-go-cms/_const/enum/user_e"
	"frozen-go-cms/common/domain"
	"frozen-go-cms/common/resource/mysql"
	"gorm.io/gorm"
)

type User struct {
	mysql.Entity
	Mobile string
	Name   string
	Gender user_e.UserGender // 0:other 1:male 2:female
	Status int
	Pwd    string
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

func UpdateUser(model *domain.Model, id mysql.ID, name string, gender user_e.UserGender) error {
	updates := map[string]interface{}{
		"name":   name,
		"gender": gender,
	}
	return model.DB().Model(User{}).Where("id = ?", id).Updates(updates).Error
}

// 获取所有用户
func ListUser(model *domain.Model) []User {
	var res []User
	if err := model.DB().Model(User{}).Find(&res).Error; err != nil {
		model.Log.Errorf("ListUser fail:%v", err)
	}
	return res
}

// 更新用户密码
func ChangeUserPwd(model *domain.Model, userId mysql.ID, pwd string) error {
	return model.DB().Model(User{}).Where("id = ?", userId).UpdateColumn("pwd", pwd).Error
}
