package router

import (
	"github.com/gin-gonic/gin"
	"hzer/internal/controller/api/captcha"
	"hzer/internal/controller/api/tests"
	"hzer/internal/controller/ws"
	"hzer/internal/middleware"
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

	//绑定静态资源
	/*if isDebug {
		r.StaticFS("/public", http.Dir("./static/public"))
	} else {
		r.StaticFS("/public", http.FS(static.Proot))
	}

	r.Use(middleware.Cors()) //取消注释此行可以开启跨域
	adminr := r.Group(configs.Data.App.AdminMain)
	{
		admin.GinApi(adminr, isDebug)
	}

	swag := r.Group("/swagger")
	{
		swag.Use(middleware.Cors())
		swag.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

		web, _ = fs.Sub(static.Js, "js")
		r.StaticFS("js", http.FS(web))
		web, _ = fs.Sub(static.Css, "css")
		r.StaticFS("css", http.FS(web))
		web, _ = fs.Sub(static.Fonts, "fonts")
		r.StaticFS("fonts", http.FS(web))

	apiR := r.Group("/api")
	{
		//api1.RouterV1(apiR)
		//验证码相关api
		captcha.GinApi(apiR)
	}*/

}
