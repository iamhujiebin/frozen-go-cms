package domain

import (
	"fmt"
	"frozen-go-cms/hilo-common/mycontext"
	"frozen-go-cms/hilo-common/resource/mysql"
	"frozen-go-cms/hilo-common/resource/redisCli"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Model struct {
	*CtxAndDb `gorm:"-"`
}

func CreateModel(ctxAndDb *CtxAndDb) *Model {
	return &Model{CtxAndDb: ctxAndDb}
}

func CreateModelContext(myContext *mycontext.MyContext) *Model {
	return &Model{
		CtxAndDb: &CtxAndDb{
			Db:        mysql.Db,
			MyContext: myContext,
			Redis:     redisCli.GetRedis(),
		},
	}
}

func CreateModelNil() *Model {
	return &Model{
		CtxAndDb: &CtxAndDb{
			Db:        mysql.Db,
			MyContext: mycontext.CreateMyContext(nil),
			Redis:     redisCli.GetRedis(),
		},
	}
}

func (m *Model) DB() *gorm.DB {
	return m.Db.WithContext(m)
}

//获取traceid
func (m *Model) GetTraceId() string {
	if traceIdTemp, ok := m.Cxt[mycontext.TRACEID]; ok {
		if traceId, ok := traceIdTemp.(string); ok {
			return traceId
		}
	}
	return ""
}

// 包装事务
// 注意:需要使用新的model
func (m *Model) Transaction(f func(*Model) error) error {
	// 公用context
	// 新的db
	txModel := CreateModelContext(m.MyContext)
	txModel.Db = m.Db.Begin().WithContext(m)
	err := f(txModel)
	if err != nil {
		txModel.Db.Rollback()
		return err
	}
	return txModel.Db.Commit().Error
}

func Persistent(db *gorm.DB, t mysql.EntityI) error {
	if t == nil {
		return nil
	}
	if t.IsLazyLoad() {
		return nil
	}
	//删除
	if t.CheckDel() {
		tx := db.Delete(t)
		if err := tx.Error; err != nil {
			return err
		}
		if tx.RowsAffected == 0 {
			return fmt.Errorf("gorm delete.RowsAffected = 0")
		}
		//增加缓存行为记录（删除）

	} else if t.GetID() == 0 {
		//新增
		if t.CheckOnDuplicateKeyUPDATE() {
			if err := db.Set("gorm:insert_option", fmt.Sprintf("ON DUPLICATE KEY UPDATE `created_time` = '%s'", time.Now())).Create(t).Error; err != nil {
				return err
			}
		} else if t.CheckOnDuplicateKeyIGNORE() {
			if err := db.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(t).Error; err != nil {
				return err
			}
		} else {
			if err := db.Create(t).Error; err != nil {
				return err
			}
		}
		//增加缓存行为记录（新增）
	} else {
		//fixme: 更新条件，目前是互斥的，应该改成且。
		//更新
		if t.CheckUpdateVersion() {
			//版本号。乐观锁更新，注意，空值不更新
			tx := db.Model(t).Where("version = ? ", t.GetUpdateVersionBefore()).Updates(t)
			if err := tx.Error; err != nil {
				return err
			}
			if tx.RowsAffected == 0 {
				return fmt.Errorf("gorm version update.RowsAffected = 0")
			}
		} else if t.CheckUpdateCondition() {
			//条件更新
			tx := db.Model(t).Where(t.GetUpdateCondition()).Updates(t)
			if err := tx.Error; err != nil {
				return err
			}
			if tx.RowsAffected == 0 {
				return fmt.Errorf("gorm condition update.RowsAffected = 0")
			}
		} else if len(t.GetOmit()) > 0 {
			if err := db.Model(t).Omit(t.GetOmit()...).Save(t).Error; err != nil {
				return err
			}
		} else {
			if err := db.Model(t).Save(t).Error; err != nil {
				return err
			}
		}
		//增加缓存行为记录（更新）
	}
	return nil
}
