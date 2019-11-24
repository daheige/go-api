package controller

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"

	"go-api/app/extensions/logger"
)

const (
	HTTP_SUCCESS_CODE = 200
	HTTP_ERROR_CODE   = 500
	API_SUCCESS_CODE  = 0
)

//用作空[]返回
type EmptyArray []struct{}

type BaseController struct{}

func (ctrl *BaseController) ajaxReturn(ctx *gin.Context, code int, message string, data interface{}) {
	// 这里不建议在这里进行select判断客户端是否超时了，如果大量的请求过来select压力太大
	// 可能会出现阻塞goroutine运行
	// if err := ctrl.ClientDisconnected(ctx); err != nil {
	// 	data = nil
	//
	// 	return
	// }

	ctx.JSON(HTTP_SUCCESS_CODE, gin.H{
		"code":     code,
		"message":  message,
		"data":     data,
		"req_time": time.Now().Unix(),
	})
}

func (ctrl *BaseController) Success(ctx *gin.Context, message string, data interface{}) {
	if len([]rune(message)) == 0 {
		message = "ok"
	}

	ctrl.ajaxReturn(ctx, API_SUCCESS_CODE, message, data)
}

//错误处理code,message
func (ctrl *BaseController) Error(ctx *gin.Context, code int, message string) {
	if code <= 0 {
		code = HTTP_ERROR_CODE
	}

	ctrl.ajaxReturn(ctx, code, message, nil)
}

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
