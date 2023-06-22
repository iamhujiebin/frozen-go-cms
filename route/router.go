package route

import (
	_ "frozen-go-cms/docs"
	"frozen-go-cms/route/album_r"
	"frozen-go-cms/route/article_r"
	"frozen-go-cms/route/channel_r"
	"frozen-go-cms/route/chatgpt_r"
	"frozen-go-cms/route/music_r"
	"frozen-go-cms/route/todo_r"
	"frozen-go-cms/route/user_r"
	"frozen-go-cms/route/ws_r"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func InitRouter() *gin.Engine {
	var r = gin.Default()
	r.Use(Cors()) // 跨域
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ws
	//r.GET("/ws", JWTApiHandle, ws_r.WsHandler)
	r.GET("/ws/:token", ws_r.WsHandler)
	r.GET("/wsTest", ws_r.WsTest)
	// http
	noLogin := r.Group("")
	noLogin.Use(ExceptionHandle, LoggerHandle)
	v1 := noLogin.Group("/v1_0")
	v1.POST("authorizations", wrapper(user_r.UserAuth))

	user := v1.Group("user")
	user.Use(JWTApiHandle)
	{
		user.GET("profile", wrapper(user_r.UserProfile))
		user.PUT("profile", wrapper(user_r.PutUserProfile))
	}
	v1.GET("channels", wrapper(channel_r.Channels))

	articles := v1.Group("articles")
	articles.Use(JWTApiHandle)
	{
		articles.POST("", wrapper(article_r.PostArticle))
		articles.PUT(":id", wrapper(article_r.PutArticle))
		articles.GET(":id", wrapper(article_r.GetArticle))
		articles.DELETE(":id", wrapper(article_r.DeleteArticle))
		articles.GET("", wrapper(article_r.PageArticle))
	}
	todo := v1.Group("todolist")
	todo.Use(JWTApiHandle)
	{
		todo.GET("", wrapper(todo_r.TodoList))
		todo.POST("", wrapper(todo_r.AddTodoList))
		todo.PUT(":id", wrapper(todo_r.MarkTodoList))
		todo.DELETE(":id", wrapper(todo_r.DelTodoList))
		todo.POST("markAll", wrapper(todo_r.MarkAllTodoList))
	}
	album := v1.Group("album")
	album.Use(JWTApiHandle)
	{
		album.GET("list", wrapper(album_r.AlbumList))
		album.POST("add", wrapper(album_r.AlbumAdd))
		album.DELETE("del/:id", wrapper(album_r.AlbumDel))
		album.GET("detail/:id", wrapper(album_r.AlbumDetail))
		album.POST("detail", wrapper(album_r.AddAlbumDetail))
	}
	chatgpt := v1.Group("chatgpt")
	chatgpt.Use(JWTApiHandle)
	{
		chatgpt.POST("process", wrapper(chatgpt_r.Process))
		chatgpt.GET("session/list", wrapper(chatgpt_r.SessionList))
		chatgpt.POST("session/add", wrapper(chatgpt_r.SessionAdd))
		chatgpt.DELETE("session/del/:id", wrapper(chatgpt_r.SessionDel))
		chatgpt.GET("session/detail/:id", wrapper(chatgpt_r.SessionDetail))
	}
	music := v1.Group("music")
	music.Use(JWTApiHandle)
	{
		music.GET("list", wrapper(music_r.MusicList))
		music.GET("search", wrapper(music_r.MusicSearch))
	}
	return r
}
