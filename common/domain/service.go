package domain

import (
	"frozen-go-cms/common/mycontext"
	"frozen-go-cms/common/resource/mysql"
	"frozen-go-cms/common/resource/redisCli"
	"runtime/debug"
)

type Service struct {
	*CtxAndDb
}

func (service *Service) getMyContext() *mycontext.MyContext {
	return service.MyContext
}

/**
 * 创建服务
 * @param
 * @return
 **/
func CreateService(myContext *mycontext.MyContext) *Service {
	if myContext == nil {
		return &Service{CtxAndDb: &CtxAndDb{
			Db:        mysql.Db,
			MyContext: mycontext.CreateMyContext(nil),
			Redis:     redisCli.GetRedis(),
		}}
	} else {
		return &Service{CtxAndDb: &CtxAndDb{
			Db:        mysql.Db,
			MyContext: myContext,
			Redis:     redisCli.GetRedis(),
		}}
	}
}

// 事务钩子回调，遇到错误，异常则回调，写service都需要钩子回调
func (service *Service) Transactional(callback func() error) error {
	//异常回调
	defer func() {
		if err := recover(); err != nil {
			service.Log.Errorf("doTransactional SYSTEM ACTION PANIC: %v, stack: %v", err, string(debug.Stack()))
			service.Db.Rollback()
			//为了防止给controller层造成数据错误，继续抛恐慌
			panic(err)
		}
	}()
	service.CtxAndDb.Db = mysql.Db.Begin()
	err := callback()
	if err != nil {
		service.Db.Rollback()
		return err
	}
	//提交
	return service.Db.Commit().Error
}
