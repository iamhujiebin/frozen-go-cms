package route

import (
	_ "frozen-go-cms/docs"
	"frozen-go-cms/route/article_r"
	"frozen-go-cms/route/channel_r"
	"frozen-go-cms/route/chatgpt_r"
	"frozen-go-cms/route/todo_r"
	"frozen-go-cms/route/user_r"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func InitRouter() *gin.Engine {
	var r = gin.Default()
	r.Use(Cors()) // 跨域
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	noLogin := r.Group("")
	noLogin.Use(ExceptionHandle, LoggerHandle)
	v1 := noLogin.Group("/v1_0")
	v1.POST("authorizations", wrapper(user_r.UserAuth))

	user := v1.Group("user")
	user.Use(JWTApiHandle)
	{
		user.GET("profile", wrapper(user_r.UserProfile))
	}
	v1.GET("channels", wrapper(channel_r.Channels))

	articles := v1.Group("mp/articles")
	articles.Use(JWTApiHandle)
	{
		articles.POST("", wrapper(article_r.PostArticle))
		articles.PUT(":id", wrapper(article_r.PutArticle))
		articles.GET(":id", wrapper(article_r.GetArticle))
		articles.DELETE(":id", wrapper(article_r.DeleteArticle))
		articles.GET("", wrapper(article_r.PageArticle))
	}
	todo := v1.Group("mp/todolist")
	todo.Use(JWTApiHandle)
	{
		todo.GET("", wrapper(todo_r.TodoList))
		todo.POST("", wrapper(todo_r.AddTodoList))
		todo.PUT(":id", wrapper(todo_r.MarkTodoList))
		todo.DELETE(":id", wrapper(todo_r.DelTodoList))
		todo.POST("markAll", wrapper(todo_r.MarkAllTodoList))
	}
	chatgpt := v1.Group("chatgpt")
	{
		chatgpt.POST("process", wrapper(chatgpt_r.Process))
	}
	return r
}
