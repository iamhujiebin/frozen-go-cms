package model

import (
	"fmt"
	"frozen-go-cms/common/resource/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

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
