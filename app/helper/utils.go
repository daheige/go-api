package helper

import (
	"regexp"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/daheige/thinkgo/monitor"
	"github.com/prometheus/client_golang/prometheus"
)

// GetDeviceByUa 根据ua获取设备名称
func GetDeviceByUa(ua string) string {
	plat := "web"
	regText := "(a|A)ndroid|dr"
	re := regexp.MustCompile(regText)
	if re.MatchString(ua) {
		plat = "android"
	} else {
		regText = "i(p|P)(hone|ad|od)|(m|M)ac"
		re = regexp.MustCompile(regText)
		if re.MatchString(ua) {
			plat = "ios"
		}
	}

	return plat
}

// Monitor metrics性能监控，gin处理器函数，包装 handler function,不侵入业务逻辑
func Monitor() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		duration := time.Since(start)
		// counter类型 metric的记录方式
		monitor.WebRequestTotal.With(prometheus.Labels{"method": ctx.Request.Method, "endpoint": ctx.Request.URL.Path}).Inc()
		// Histogram类型 meric的记录方式
		monitor.WebRequestDuration.With(prometheus.Labels{"method": ctx.Request.Method, "endpoint": ctx.Request.URL.Path}).Observe(duration.Seconds())
	}

}
