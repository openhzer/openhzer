package captcha

import "github.com/gin-gonic/gin"

func GinApi(g *gin.RouterGroup) {
	//r := g.Group("/captcha")
	g.GET("/captcha", getCaptcha)
	//g.POST("/captcha", verifyCaptcha)
}
