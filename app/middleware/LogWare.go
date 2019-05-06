package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type LogWare struct{}

func (ware *LogWare) Access() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		t := time.Now()
		uri := ctx.Request.RequestURI
		log.Println("request uri: ", uri)

		//测试上下文设置，将id添加到上下文中
		ctx.Set("current_uid", ctx.DefaultQuery("id", "123"))
		ctx.Next()

		log.Println("request end")
		log.Println("response cost time: ", time.Now().Sub(t))
		return
	}
}
