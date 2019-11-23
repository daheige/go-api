package Logger

import (
	"go-api/app/helper"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/daheige/thinkgo/logger"
)

/**
 * {
    "level":"info",
    "time_local":"2019-06-29T12:00:22.936+0800",
    "line":"/web/go/go-api/app/extensions/Logger/log.go:59",
    "msg":"exec start",
    "ip":"::1",
    "ua":"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.80 Safari/537.36",
    "trace_line":30,
    "log_id":"5197eeb0b7ec2847fe5fb412cd8f9ebc",
    "options":null,
    "host":"[::1]:36850",
    "plat":"web",
    "method":"GET",
    "trace_file":"/web/go/go-api/app/middleware/LogWare.go",
    "tag":"_",
    "request_uri":"/"
}
*/

func writeLog(ctx *gin.Context, levelName string, message string, context map[string]interface{}) {
	tag := strings.Replace(ctx.Request.RequestURI, "/", "_", -1)
	ua := ctx.GetHeader("User-Agent")

	//log.Println("context: ", context)

	//函数调用
	_, file, line, _ := runtime.Caller(2)
	logInfo := map[string]interface{}{
		"tag":         tag,
		"request_uri": ctx.Request.RequestURI,
		"log_id":      ctx.GetString("log_id"),
		"options":     context,
		"host":        ctx.Request.RemoteAddr,
		"ip":          ctx.ClientIP(),
		"ua":          ua,
		"plat":        helper.GetDeviceByUa(ua), //当前设备匹配
		"method":      ctx.Request.Method,
		"trace_line":  line,
		"trace_file":  file,
	}

	switch levelName {
	case "info":
		logger.Info(message, logInfo)
	case "debug":
		logger.Debug(message, logInfo)
	case "warn":
		logger.Warn(message, logInfo)
	case "error":
		logger.Error(message, logInfo)
	case "emergency":
		logger.DPanic(message, logInfo)
	default:
		logger.Info(message, logInfo)
	}
}

func Info(ctx *gin.Context, message string, context map[string]interface{}) {
	writeLog(ctx, "info", message, context)
}

func Debug(ctx *gin.Context, message string, context map[string]interface{}) {
	writeLog(ctx, "debug", message, context)
}

func Warn(ctx *gin.Context, message string, context map[string]interface{}) {
	writeLog(ctx, "warn", message, context)
}

func Error(ctx *gin.Context, message string, context map[string]interface{}) {
	writeLog(ctx, "error", message, context)
}

//致命错误或panic捕获
func Emergency(ctx *gin.Context, message string, context map[string]interface{}) {
	writeLog(ctx, "emergency", message, context)
}

//异常捕获处理
func Recover(c interface{}) {
	defer func() {
		if err := recover(); err != nil {
			if ctx, ok := c.(*gin.Context); ok {
				Emergency(ctx, "exec panic", map[string]interface{}{
					"error":       err,
					"error_trace": string(debug.Stack()),
				})

				//响应状态
				ctx.AbortWithStatusJSON(500, gin.H{
					"code":    500,
					"message": "server error",
				})

				return
			}

			logger.DPanic("exec panic", map[string]interface{}{
				"error":       err,
				"error_trace": string(debug.Stack()),
			})
		}
	}()

}
