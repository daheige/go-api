package routes

import (
	"go-api/app/controller"
	"go-api/app/helper"
	"go-api/app/middleware"

	"github.com/gin-gonic/gin"
)

func WebRoute(router *gin.Engine) {
	//访问日志中间件处理
	logWare := &middleware.LogWare{}

	//对所有的请求进行性能监控，一般来说生产环境，可以对指定的接口做性能监控
	router.Use(logWare.Access(), logWare.Recover(), helper.Monitor())
	//router.Use(logWare.Access(), logWare.Recover())

	router.NoRoute(middleware.NotFoundHandler())

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

	//对个别接口进行性能监控
	//router.GET("/index", helper.Monitor(), homeCtrl.Index)
	router.GET("/index", homeCtrl.Index)

	router.GET("/test", homeCtrl.Test)

	//定义api前缀分组
	v1 := router.Group("/v1")
	// http://localhost:1338/v1/info/123
	v1.GET("/info/:id", homeCtrl.Info)

	//http://localhost:1338/v1/data?id=456
	v1.GET("/data", homeCtrl.GetData)

	v1.GET("/set-data", homeCtrl.SetData)

	v1.POST("/data", homeCtrl.PostData)

	router.GET("/long-async", homeCtrl.LongAsync)
}
