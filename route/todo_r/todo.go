package todo_r

import (
	"frozen-go-cms/domain/model/todo_m"
	"frozen-go-cms/hilo-common/domain"
	"frozen-go-cms/hilo-common/mycontext"
	"frozen-go-cms/myerr/bizerr"
	"frozen-go-cms/req"
	"frozen-go-cms/resp"
	"github.com/gin-gonic/gin"
	"sort"
	"strconv"
)

type CvTodoList struct {
	Id     uint64 `json:"id"`
	Name   string `json:"name"`
	IsDone bool   `json:"isDone"`
}

// @Tags Todo模块
// @Summary 列表
// @Param Authorization header string true "token"
// @Success 200 {object} []CvTodoList
// @Router /v1_0/mp/todolist [get]
func TodoList(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	userId, err := req.GetUserId(c)
	if err != nil {
		return myCtx, err
	}
	model := domain.CreateModelContext(myCtx)
	todos := todo_m.GetTodos(model, userId)
	var todolist []CvTodoList
	for _, todo := range todos {
		todolist = append(todolist, CvTodoList{
			Id:     todo.ID,
			Name:   todo.Name,
			IsDone: todo.IsDone,
		})
	}
	// isDone的放在后面
	sort.SliceStable(todolist, func(i, j int) bool {
		if !todolist[i].IsDone && todolist[j].IsDone {
			return true
		}
		return false
	})
	resp.ResponseOk(c, todolist)
	return myCtx, nil
}

type AddTodoListReq struct {
	Name string `json:"name" binding:"required"`
}

// @Tags Todo模块
// @Summary 添加
// @Param Authorization header string true "token"
// @Param AddTodoListReq body AddTodoListReq true "请求体"
// @Success 200
// @Router /v1_0/mp/todolist [post]
func AddTodoList(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	var param AddTodoListReq
	if err := c.ShouldBind(&param); err != nil {
		return myCtx, err
	}
	userId, err := req.GetUserId(c)
	if err != nil {
		return myCtx, err
	}
	model := domain.CreateModelContext(myCtx)
	var id uint64
	if id, err = todo_m.AddTodo(model, todo_m.Todo{UserId: userId, Name: param.Name, IsDone: false}); err != nil {
		return myCtx, err

	}
	resp.ResponseOk(c, id)
	return myCtx, nil
}

type MarkTodoListReq struct {
	IsDone bool `json:"isDone"`
}

// @Tags Todo模块
// @Summary 标记
// @Param Authorization header string true "token"
// @Param id path integer true "id"
// @Param MarkTodoListReq body MarkTodoListReq true "请求体"
// @Success 200
// @Router /v1_0/mp/todolist/:id [put]
func MarkTodoList(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	_id := c.Param("id")
	id, _ := strconv.ParseUint(_id, 10, 64)
	if id <= 0 {
		return myCtx, bizerr.ParaMissing
	}
	var param MarkTodoListReq
	if err := c.ShouldBind(&param); err != nil {
		return myCtx, err
	}
	model := domain.CreateModelContext(myCtx)
	if err := todo_m.MarkTodo(model, id, param.IsDone); err != nil {
		return myCtx, err
	}
	resp.ResponseOk(c, "")
	return myCtx, nil
}

type MarkAllTodoListReq struct {
	IsDone bool `json:"isDone"`
}

// @Tags Todo模块
// @Summary 标记全部
// @Param Authorization header string true "token"
// @Param MarkAllTodoListReq body MarkAllTodoListReq true "请求体"
// @Success 200
// @Router /v1_0/mp/todolist/markAll [post]
func MarkAllTodoList(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	var param MarkAllTodoListReq
	if err := c.ShouldBind(&param); err != nil {
		return myCtx, err
	}
	userId, err := req.GetUserId(c)
	if err != nil {
		return myCtx, err
	}
	model := domain.CreateModelContext(myCtx)
	if err := todo_m.MarkTodoAll(model, userId, param.IsDone); err != nil {
		return myCtx, err
	}
	resp.ResponseOk(c, "")
	return myCtx, nil
}

// @Tags Todo模块
// @Summary 删除
// @Param id path integer true "id"
// @Success 200
// @Router /v1_0/mp/todolist/:id [delete]
func DelTodoList(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	_id := c.Param("id")
	id, _ := strconv.ParseUint(_id, 10, 64)
	if id <= 0 {
		return myCtx, bizerr.ParaMissing
	}
	model := domain.CreateModelContext(myCtx)
	if err := todo_m.DelTodo(model, id); err != nil {
		return myCtx, err
	}
	resp.ResponseOk(c, "")
	return myCtx, nil
}
