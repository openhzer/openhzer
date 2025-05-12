package router

import (
	"github.com/gin-gonic/gin"
	ginStatic "github.com/soulteary/gin-static"
	"hzer/internal/controller/api/captcha"
	"hzer/internal/controller/api/tests"
	"hzer/internal/controller/ws"
	"hzer/internal/middleware"
	"hzer/static"
)

func NewHTTPRouter(r *gin.Engine) {
	//isDebug := os.Getenv("GIN_MODE") == "debug"
	rootRouter := r.Group("/")
	apiRouter := r.Group("/api")

	/*swag := r.Group("/swagger")
	{
		swag.Use(middleware.Cors())
		swag.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}*/

	rootRouter.Use(middleware.Cors())
	apiRouter.Use(middleware.Cors())

	captcha.GinApi(apiRouter)
	ws.GinApi(rootRouter)
	tests.GinApi(apiRouter)

	// 静态资源
	r.Use(ginStatic.ServeEmbed("public", static.Public))

}
