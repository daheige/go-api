package logger

import (
	"context"
	"go-api/app/helper"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/daheige/thinkgo/gutils"
	"github.com/daheige/thinkgo/logger"
)

/**
{
    "level":"info",
    "time_local":"2019-11-24T19:36:02.838+0800",
    "msg":"exec start",
    "request_method":"GET",
    "request_uri":"/v1/hello",
    "log_id":"c34a00fa-0906-2fff-b5e8-61cfcc2a19ff",
    "ua":"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.80 Safari/537.36",
    "plat":"web",
    "trace_line":41,
    "trace_file":"/web/go/go-api/app/middleware/log.go",
    "tag":"v1_hello",
    "options":null,
    "ip":"127.0.0.1"
}
*/

func writeLog(ctx context.Context, levelName string, message string, options map[string]interface{}) {
	reqUri := getStringByCtx(ctx, "request_uri")
	tag := strings.Replace(reqUri, "/", "_", -1)
	tag = strings.Replace(tag, ".", "_", -1)
	tag = strings.TrimLeft(tag, "_")

	logId := getStringByCtx(ctx, "log_id")
	if logId == "" {
		logId = gutils.RndUuid()
	}

	ua := getStringByCtx(ctx, "user_agent")

	//函数调用
	_, file, line, _ := runtime.Caller(2)
	logInfo := map[string]interface{}{
		"tag":            tag,
		"request_uri":    reqUri,
		"log_id":         logId,
		"options":        options,
		"ip":             getStringByCtx(ctx, "client_ip"),
		"ua":             ua,
		"plat":           helper.GetDeviceByUa(ua), //当前设备匹配
		"request_method": getStringByCtx(ctx, "request_method"),
		"trace_line":     line,
		"trace_file":     file,
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

func getStringByCtx(ctx context.Context, key string) string {
	return helper.GetStringByCtx(ctx, key)
}

func Info(ctx context.Context, message string, context map[string]interface{}) {
	writeLog(ctx, "info", message, context)
}

func Debug(ctx context.Context, message string, context map[string]interface{}) {
	writeLog(ctx, "debug", message, context)
}

func Warn(ctx context.Context, message string, context map[string]interface{}) {
	writeLog(ctx, "warn", message, context)
}

func Error(ctx context.Context, message string, context map[string]interface{}) {
	writeLog(ctx, "error", message, context)
}

//致命错误或panic捕获
func Emergency(ctx context.Context, message string, context map[string]interface{}) {
	writeLog(ctx, "emergency", message, context)
}

//异常捕获处理
func Recover(c interface{}) {
	defer func() {
		if err := recover(); err != nil {
			if ctx, ok := c.(context.Context); ok {
				Emergency(ctx, "exec panic", map[string]interface{}{
					"error":       err,
					"error_trace": string(debug.Stack()),
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
