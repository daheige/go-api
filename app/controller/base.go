package controller

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/daheige/go-api/app/extensions/logger"
)

const (
	// HttpSuccessCode http success code
	HttpSuccessCode = 200

	// HttpErrorCode http error code
	HttpErrorCode = 500

	// ApiSuccessCode api success logic code
	ApiSuccessCode = 0
)

// EmptyArray用作空[]返回
type EmptyArray []struct{}

// BaseController base ctrl.
type BaseController struct{}

func (ctrl *BaseController) ajaxReturn(ctx *gin.Context, code int, message string, data interface{}) {
	// 这里不建议在这里进行select判断客户端是否超时了，如果大量的请求过来select压力太大
	// 可能会出现阻塞goroutine运行
	// if err := ctrl.ClientDisconnected(ctx); err != nil {
	// 	data = nil
	//
	// 	return
	// }

	ctx.JSON(HttpSuccessCode, gin.H{
		"code":     code,
		"message":  message,
		"data":     data,
		"req_time": time.Now().Unix(),
	})
}

// Success success
func (ctrl *BaseController) Success(ctx *gin.Context, message string, data interface{}) {
	if len([]rune(message)) == 0 {
		message = "ok"
	}

	ctrl.ajaxReturn(ctx, ApiSuccessCode, message, data)
}

// Error 错误处理code,message
func (ctrl *BaseController) Error(ctx *gin.Context, code int, message string) {
	if code <= 0 {
		code = HttpErrorCode
	}

	ctrl.ajaxReturn(ctx, code, message, nil)
}

// ClientDisconnected client disconnect.
func (ctrl *BaseController) ClientDisconnected(c *gin.Context) error {

	// 标准上下文
	ctx := c.Request.Context()

	select {
	// if the context is done it timed out or was cancelled
	// so don't return anything
	case <-ctx.Done():
		logger.Error(c, "client disconnected", map[string]interface{}{
			"trace_error": ctx.Err().Error(),
		})

		return errors.New("client disconnected")
	default:
	}

	return nil
}
