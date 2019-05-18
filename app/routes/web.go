package routes

import (
	"go-api/app/controller"
	"go-api/app/middleware"

	"github.com/gin-gonic/gin"
)

func WebRoute(router *gin.Engine) {
	//访问日志中间件处理
	logWare := &middleware.LogWare{}
	router.Use(logWare.Access(), logWare.Recover())

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"code":    0,
			"message": "ok",
		})
	})

	router.GET("/check", func(ctx *gin.Context) {
		ctx.String(200, `{"alive": true}`)
	})

	homeCtrl := &controller.HomeController{}
	router.GET("/index", homeCtrl.Index)
	router.GET("/test", homeCtrl.Test)

	//定义api前缀分组
	v1 := router.Group("/v1")
	v1.GET("/info", homeCtrl.Info)

	//http://localhost:1338/v1/data?id=456
	v1.GET("/data", homeCtrl.GetData)

	v1.GET("/set-data", homeCtrl.SetData)
}
