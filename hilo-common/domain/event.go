package domain

import "frozen-go-cms/hilo-common/utils"

//异步执行的接口-老方式-废弃
//type AsyncEvent interface {
//	AsyncDo(model *Model, eventData interface{}, n int) error
//	AsyncSize() int
//	AsyncNoTxDo(model *Model, eventData interface{}, n int) error
//	AsyncNoTxSize() int
//}

// 程序内部事件
type EventBase struct {
	//同步执行
	syncList []func(model *Model, event interface{}) error
	//异步执行
	asyncList []func(model *Model, event interface{}) error
}

// 添加同步事件
func AddEventSync(event *EventBase, callback func(model *Model, event interface{}) error) {
	event.syncList = append(event.syncList, callback)
}

// 添加异步事件
func AddEventAsync(event *EventBase, callback func(model *Model, event interface{}) error) {
	event.asyncList = append(event.asyncList, callback)
}

// 发布事件
func PublishEvent(event *EventBase, model *Model, data interface{}) error {
	// 执行同步的领域事件
	for _, callback := range event.syncList {
		if err := callback(model, data); err != nil {
			return err
		}
	}
	// 执行异步的领域事件
	if len(event.asyncList) > 0 {
		go func() {
			defer utils.CheckGoPanic()
			for _, callback := range event.asyncList {
				// 异步事件需要用新model,主要是db
				var newModel = CreateModelContext(model.MyContext)
				if err := callback(newModel, data); err != nil {
					model.Log.Errorf("aysnc fail:%v", err)
				}
			}
		}()
	}
	return nil
}
