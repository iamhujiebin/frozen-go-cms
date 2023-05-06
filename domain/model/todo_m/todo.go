package todo_m

import (
	"git.hilo.cn/hilo-common/domain"
	"git.hilo.cn/hilo-common/resource/mysql"
)

type Todo struct {
	mysql.Entity
	UserId mysql.ID
	Name   string
	IsDone bool
}

func AddTodo(model *domain.Model, todo Todo) (uint64, error) {
	err := model.DB().Create(&todo).Error
	return todo.ID, err
}

func DelTodo(model *domain.Model, id mysql.ID) error {
	return model.DB().Model(Todo{}).Where("id = ?", id).Delete(&Todo{}).Error
}

func MarkTodo(model *domain.Model, id mysql.ID, isDone bool) error {
	return model.DB().Model(Todo{}).Where("id = ?", id).Update("is_done", isDone).Error
}

func MarkTodoAll(model *domain.Model, userId mysql.ID, isDone bool) error {
	return model.DB().Model(Todo{}).Where("user_id = ?", userId).Update("is_done", isDone).Error
}

func GetTodos(model *domain.Model, userId mysql.ID) []Todo {
	var todos []Todo
	if err := model.DB().Model(Todo{}).Where("user_id = ?", userId).Order("id DESC").Find(&todos).Error; err != nil {
		model.Log.Errorf("GetTodos fail:%v", err)
	}
	return todos
}
