package mysql

import "time"

type EntityI interface {
	GetID() ID
	//用于判断数据是否进行持久化
	IsLazyLoad() bool
	//默认值为false true:代表要移除数据
	CheckDel() bool
	//检查是否唯一键冲突，依旧更新
	CheckOnDuplicateKeyUPDATE() bool
	//检查是否唯一键冲突，则不插入
	CheckOnDuplicateKeyIGNORE() bool
	//更新乐观锁 默认值为false true:乐观锁更新
	CheckUpdateVersion() bool
	//更新条件.
	CheckUpdateCondition() bool
	//获取版本号
	GetUpdateVersionBefore() uint
	//更新情况
	GetUpdateCondition() string
	//save 动作排除字段
	GetOmit() []string
}
type Entity struct {
	ID                   ID                     `gorm:"primary_key"`
	CreatedTime          time.Time              `gorm:"->"`
	UpdatedTime          time.Time              `gorm:"->"`
	lazyLoad             bool                   `gorm:"-"`
	del                  bool                   `gorm:"-"`
	onDuplicateKeyUPDATE bool                   `gorm:"-"`
	onDuplicateKeyIGNORE bool                   `gorm:"-"`
	updateVersionFlag    bool                   `gorm:"-"`
	updateVersionBefore  uint                   `gorm:"-"`
	updateCondition      string                 `gorm:"-"`
	omit                 []string               `gorm:"-"` //更新排除
	updateColumns        map[string]interface{} `gorm:"-"` //更新字段
}

func (t *Entity) GetID() ID {
	return t.ID
}

func (t *Entity) IsLazyLoad() bool {
	return t.lazyLoad
}

func (t *Entity) SetLasyLoad() {
	t.lazyLoad = true
}

func (t *Entity) SetDel() {
	t.del = true
}

func (t *Entity) CheckDel() bool {
	return t.del
}

func (t *Entity) SetOnDuplicateKeyUPDATE() {
	t.onDuplicateKeyUPDATE = true
}

func (t *Entity) SetOnDuplicateKeyIGNORE() {
	t.onDuplicateKeyIGNORE = true
}

func (t *Entity) CheckOnDuplicateKeyUPDATE() bool {
	return t.onDuplicateKeyUPDATE
}

func (t *Entity) CheckOnDuplicateKeyIGNORE() bool {
	return t.onDuplicateKeyIGNORE
}

func (t *Entity) SetCheckUpdateVersionBefore(versionBefore uint) {
	t.updateVersionBefore = versionBefore
	t.updateVersionFlag = true
}

func (t *Entity) SetCheckUpdateCondition(condition string) {
	t.updateCondition = condition
}

func (t *Entity) CheckUpdateVersion() bool {
	return t.updateVersionFlag
}

func (t *Entity) CheckUpdateCondition() bool {
	return t.updateCondition != ""
}

func (t *Entity) GetUpdateCondition() string {
	return t.updateCondition
}

func (t *Entity) GetUpdateVersionBefore() uint {
	return t.updateVersionBefore
}

func (t *Entity) GetOmit() []string {
	return t.omit
}

func (t *Entity) SetOmit(omit []string) {
	t.omit = omit
}

func (t *Entity) SetUpdateColumns(updateColumns map[string]interface{}) {
	t.updateColumns = updateColumns
}

func (t *Entity) GetUpdateColumns() map[string]interface{} {
	return t.updateColumns
}
