package routes

import (
	"time"

	"github.com/daheige/go-api/app/controller"
	"github.com/daheige/go-api/app/helper"
	"github.com/daheige/go-api/app/middleware"

	"github.com/gin-gonic/gin"
)

// WebRoute web router.
func WebRoute(router *gin.Engine) {
	// 心跳检测
	router.GET("/check", func(ctx *gin.Context) {
		ctx.String(200, `{"alive": true}`)
	})

	//访问日志中间件和recover捕获
	logWare := &middleware.LogWare{}
	router.Use(logWare.Access(), logWare.Recover())

	// 服务超时设置 3s超时
	router.Use(middleware.TimeoutHandler(3 * time.Second))

	// prometheus监控
	// 对所有的请求进行性能监控，一般来说生产环境，可以对指定的接口做性能监控
	router.Use(helper.Monitor())

	router.NoRoute(middleware.NotFoundHandler())

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"code":    0,
			"message": "ok",
		})
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

	// http://localhost:1338/v1/get-data?name=daheige
	v1.GET("/get-data", homeCtrl.GetData)

	v1.GET("/get-user", homeCtrl.GetUser)

	v1.POST("/post-data", homeCtrl.PostData)

	router.GET("/long-async", homeCtrl.LongAsync)

	//压力测试/api/info接口
	router.GET("/api/info", homeCtrl.GetInfo)

	// 压力测试map gc
	indexCtrl := &controller.IndexController{}
	v1.GET("/hello", indexCtrl.Hello)

	//模拟panic操作
	v1.GET("/test-panic", homeCtrl.Test)

	// 验证gin param validate 参数检验功能
	// http://localhost:1338/v1/person-info?id=12&address=fefefe
	v1.GET("/person-info", homeCtrl.ValidData)
}
