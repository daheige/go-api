package Logger

import (
	"go-api/app/helper"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/daheige/thinkgo/common"
)

const logTmFmtWithMS = "2006-01-02 15:04:05.999"

func writeLog(ctx *gin.Context, levelName string, message interface{}, context map[string]interface{}) {
	tag := strings.Replace(ctx.Request.RequestURI, "/", "_", -1)
	ua := ctx.GetHeader("User-Agent")

	//log.Println("context: ", context)

	//函数调用
	_, file, line, _ := runtime.Caller(2)
	logInfo := map[string]interface{}{
		"tag":          tag,
		"request_uri":  ctx.Request.RequestURI,
		"log_id":       ctx.GetString("log_id"),
		"local_time":   time.Now().Format(logTmFmtWithMS),
		"detail":       context,
		"host":         ctx.Request.RemoteAddr,
		"ua":           ua,
		"plat":         helper.GetDeviceByUa(ua), //当前设备匹配
		"method":       ctx.Request.Method,
		"current_line": line,
		"current_file": file,
	}

	switch levelName {
	case "info":
		common.InfoLog(message, logInfo)
	case "debug":
		common.DebugLog(message, logInfo)
	case "warn":
		common.WarnLog(message, logInfo)
	case "error":
		common.ErrorLog(message, logInfo)
	case "emergency":
		common.EmergLog(message, logInfo)
	}
}

func Info(ctx *gin.Context, message interface{}, context map[string]interface{}) {
	writeLog(ctx, "info", message, context)
}

func Debug(ctx *gin.Context, message interface{}, context map[string]interface{}) {
	writeLog(ctx, "debug", message, context)
}

func Warn(ctx *gin.Context, message interface{}, context map[string]interface{}) {
	writeLog(ctx, "warn", message, context)
}

func Error(ctx *gin.Context, message interface{}, context map[string]interface{}) {
	writeLog(ctx, "error", message, context)
}

//致命错误或panic捕获
func Emergency(ctx *gin.Context, message interface{}, context map[string]interface{}) {
	writeLog(ctx, "emergency", message, context)
}
