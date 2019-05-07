package middleware

import (
	"go-api/app/extensions/Logger"
	"log"
	"time"

	"github.com/daheige/thinkgo/common"

	"github.com/gin-gonic/gin"
)

type LogWare struct{}

func (ware *LogWare) Access() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		t := time.Now()
		uri := ctx.Request.RequestURI

		log.Println("request before")
		log.Println("request uri: ", uri)

		//如果采用了nginx x-request-id功能，可以获得x-request-id
		logId := ctx.GetString("X-Request-Id")
		if logId == "" {
			logId = common.RndUuidMd5() //日志id
		}

		ctx.Set("log_id", logId)
		Logger.Info(ctx, "exec start", nil)

		ctx.Next()

		log.Println("request end")
		//请求结束记录日志
		c := map[string]interface{}{
			"exec_time": time.Now().Sub(t).Seconds(),
		}

		if code := ctx.Writer.Status(); code != 200 {
			c["response_code"] = code
		}

		Logger.Info(ctx, "exec end", c)
	}
}

//请求处理中遇到异常或panic捕获
func (ware *LogWare) Recover() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("error:%v", err)
				Logger.Emergency(ctx, "exec panic", map[string]interface{}{
					"trace_error": err,
					"trace_info":  string(common.CatchStack()),
				})

				//响应状态
				ctx.AbortWithStatusJSON(500, gin.H{
					"code":    500,
					"message": "server error",
				})

				return
			}
		}()

		ctx.Next()
	}
}
