package route

import (
	_ "frozen-go-cms/docs"
	"frozen-go-cms/route/user_r"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func InitRouter() *gin.Engine {
	var r = gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	needLogin := r.Group("")
	needLogin.Use(ExceptionHandle, LoggerHandle)
	v1 := needLogin.Group("/v1")
	user := v1.Group("/user")
	{
		user.GET("/bag/:resType", wrapper(user_r.UserBag))
	}
	return r
}
