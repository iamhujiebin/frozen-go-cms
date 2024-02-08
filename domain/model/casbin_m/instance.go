package casbin_m

import (
	"errors"
	"frozen-go-cms/common/resource/config"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/util"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cast"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"net/url"
	"strings"
)

var enforcer *casbin.Enforcer

func init() {
	mysqlConfigData := config.GetConfigMysql()
	options := "?charset=utf8mb4&parseTime=True&loc=Local&time_zone=" + url.QueryEscape("'+8:00'")
	dsn := "" + mysqlConfigData.MYSQL_USERNAME + ":" + mysqlConfigData.MYSQL_PASSWORD + "@(" + mysqlConfigData.MYSQL_HOST + ")/" + mysqlConfigData.MYSQL_DB + options
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		panic(err)
	}
	enforcer, err = casbin.NewEnforcer("./rbac_model.conf", adapter)
	if err != nil {
		panic(err)
	}
	enforcer.AddFunction("ParamsMatch", ParamsMatchFunc)
	_ = enforcer.LoadPolicy()
}

func Casbin() *casbin.Enforcer {
	return enforcer
}

type CasbinModel struct {
	PType  string `json:"p_type" gorm:"column:p_type" description:"策略类型"`
	UserId string `json:"role_id" gorm:"column:v0" description:"角色ID"`
	Path   string `json:"path" gorm:"column:v1" description:"api路径"`
	Method string `json:"method" gorm:"column:v2" description:"访问方法"`
}

func (c *CasbinModel) TableName() string {
	return "casbin_rule"
}

func (c *CasbinModel) Create() error {
	e := Casbin()
	if success, err := e.AddPolicy(c.UserId, c.Path, c.Method); err != nil {
		return err
	} else if success == false {
		return errors.New("存在相同的API，添加失败")
	}
	return nil
}

//func (c *CasbinModel) Update(db *gorm.DB, values interface{}) error {
//	if err := db.Model(c).Where("v1 = ? AND v2 = ?", c.Path, c.Method).Updates(values).Error; err != nil {
//		return err
//	}
//	return nil
//}

func (c *CasbinModel) List() [][]string {
	e := Casbin()
	policy := e.GetFilteredPolicy(0, cast.ToString(c.UserId))
	return policy
}

// @function: ClearCasbin
// @description: 清除匹配的权限
// @param: v int, p ...string
// @return: bool
func (c *CasbinModel) ClearCasbin() error {
	e := Casbin()
	_, err := e.RemoveFilteredPolicy(0, cast.ToString(c.UserId), c.Path, c.Method)
	return err
}

// @function: ParamsMatch
// @description: 自定义规则函数
// @param: fullNameKey1 string, key2 string
// @return: bool
func ParamsMatch(fullNameKey1 string, key2 string) bool {
	key1 := strings.Split(fullNameKey1, "?")[0]
	// 剥离路径后再使用casbin的keyMatch2
	return util.KeyMatch2(key1, key2)
}

// @function: ParamsMatchFunc
// @description: 自定义规则函数
// @param: args ...interface{}
// @return: interface{}, error
func ParamsMatchFunc(args ...interface{}) (interface{}, error) {
	name1 := args[0].(string)
	name2 := args[1].(string)

	return ParamsMatch(name1, name2), nil
}
