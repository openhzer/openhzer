package tests

import (
	"github.com/gin-gonic/gin"
	"hzer/internal/response"
	"io"
	"os"
)

func GinApi(apiRouter *gin.RouterGroup) {
	r := apiRouter.Group("/test")
	r.GET("/mk", func(context *gin.Context) {
		file, _ := os.Open("README.md")
		defer file.Close()
		by, _ := io.ReadAll(file)
		response.SuccessJson(context, "ok", gin.H{"md": string(by)})
	})
	redisRouter := r.Group("/redis")
	{
		redisRouter.GET("json", redisJson)
	}
	/*	mongoRouter := r.Group("/mongo")
		{
			mongoRouter.GET("/insert_log", insertLog)
		}*/

}
