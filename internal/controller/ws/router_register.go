package ws

import "github.com/gin-gonic/gin"

func GinApi(apiRouter *gin.RouterGroup) {
	wsRouter := apiRouter.Group("/ws")
	wsRouter.GET("/test", testws)
	go Write()
}
