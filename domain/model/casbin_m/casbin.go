package casbin_m

import (
	"frozen-go-cms/common/domain"
	"github.com/spf13/cast"
)

func CasbinCreate(model *domain.Model, uid uint64, path, method string) error {
	cm := CasbinModel{
		PType:  "p",
		UserId: cast.ToString(uid),
		Path:   path,
		Method: method,
	}
	return cm.Create()
}

func CasbinRemove(model *domain.Model, uid uint64, path, method string) error {
	cm := CasbinModel{
		PType:  "p",
		UserId: cast.ToString(uid),
		Path:   path,
		Method: method,
	}
	return cm.ClearCasbin()
}

func CasbinList(model *domain.Model, uid uint64) [][]string {
	cm := CasbinModel{UserId: cast.ToString(uid)}
	return cm.List()
}
