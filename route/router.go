package route

import (
	_ "frozen-go-cms/docs"
	"frozen-go-cms/route/ai_r"
	"frozen-go-cms/route/album_r"
	"frozen-go-cms/route/article_r"
	"frozen-go-cms/route/channel_r"
	"frozen-go-cms/route/chatgpt_r"
	"frozen-go-cms/route/music_r"
	"frozen-go-cms/route/product_price_r"
	"frozen-go-cms/route/todo_r"
	"frozen-go-cms/route/user_r"
	"frozen-go-cms/route/vap_r"
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
		music.GET("down", wrapper(music_r.MusicDown))
		music.DELETE(":id", wrapper(music_r.MusicDel))
		music.GET("/author/search", wrapper(music_r.MusicAuthorSearch))
		music.GET("/author/list", wrapper(music_r.MusicAuthorList))
		music.POST("/author/down", wrapper(music_r.MusicAuthorDown))
	}
	vap := v1.Group("vap")
	vap.Use(JWTApiHandle)
	{
		vap.POST("vapc", wrapper(vap_r.VapVapc))
	}
	ai := v1.Group("ai")
	ai.Use(JWTApiHandle)
	{
		ai.GET("prompts", wrapper(ai_r.Prompts))
		ai.GET("images", wrapper(ai_r.Images))
		ai.POST("images", wrapper(ai_r.GenImages))
	}
	productPrice := v1.Group("productPrice")
	{
		// 系统配置
		productPrice.GET("/system/config", wrapper(product_price_r.SystemConfigGet))
		productPrice.PUT("/system/config", wrapper(product_price_r.SystemConfigPut))
		// 印刷价格
		productPrice.GET("/color", wrapper(product_price_r.ColorPriceGet))
		productPrice.PUT("/color/:id", wrapper(product_price_r.ColorPricePut))
		productPrice.POST("/color", wrapper(product_price_r.ColorPricePost))
		productPrice.DELETE("/color/:id", wrapper(product_price_r.ColorPriceDelete))
		// 工艺价格
		productPrice.GET("/craft", wrapper(product_price_r.CraftPriceGet))
		productPrice.PUT("/craft/:id", wrapper(product_price_r.CraftPricePut))
		productPrice.POST("/craft", wrapper(product_price_r.CraftPricePost))
		productPrice.DELETE("/craft/:id", wrapper(product_price_r.CraftPriceDelete))
		// 材料价格
		productPrice.GET("/material", wrapper(product_price_r.MaterialPriceGet))
		productPrice.PUT("/material/:id", wrapper(product_price_r.MaterialPricePut))
		productPrice.POST("/material", wrapper(product_price_r.MaterialPricePost))
		productPrice.DELETE("/material/:id", wrapper(product_price_r.MaterialPriceDelete))
		// 规格尺寸
		productPrice.GET("/size", wrapper(product_price_r.SizeConfigGet))
		productPrice.PUT("/size/:id", wrapper(product_price_r.SizeConfigPut))
		productPrice.POST("/size", wrapper(product_price_r.SizeConfigPost))
		productPrice.DELETE("/size/:id", wrapper(product_price_r.SizeConfigDelete))
	}
	return r
}
