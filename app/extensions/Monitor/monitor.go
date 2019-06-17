package Monitor

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/prometheus/client_golang/prometheus"
)

// 初始化 web_reqeust_total， counter类型指标， 表示接收http请求总次数
var WebRequestTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "web_reqeust_total",
		Help: "Number of hello requests in total",
	},
	[]string{"method", "endpoint"}, //设置两个标签 请求方法和 路径 对请求总次数在两个
)

// web_request_duration_seconds，Histogram类型指标，bucket代表duration的分布区间
var WebRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "web_request_duration_seconds",
		Help:    "web request duration distribution",
		Buckets: []float64{0.1, 0.3, 0.5, 0.7, 0.9, 1},
	},
	[]string{"method", "endpoint"},
)

var CpuTemp = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "cpu_temperature_celsius",
	Help: "Current temperature of the CPU.",
})

var HdFailures = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "hd_errors_total",
		Help: "Number of hard-disk errors.",
	},
	[]string{"device"},
)

// 对于http原始的处理器函数，包装 handler function,不侵入业务逻辑
func MonitorHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		h(w, r)

		duration := time.Since(start)
		// counter类型 metric的记录方式
		WebRequestTotal.With(prometheus.Labels{"method": r.Method, "endpoint": r.URL.Path}).Inc()
		// Histogram类型 meric的记录方式
		WebRequestDuration.With(prometheus.Labels{"method": r.Method, "endpoint": r.URL.Path}).Observe(duration.Seconds())
	}
}

// 对gin处理器函数，包装 handler function,不侵入业务逻辑
func Monitor() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		duration := time.Since(start)
		// counter类型 metric的记录方式
		WebRequestTotal.With(prometheus.Labels{"method": ctx.Request.Method, "endpoint": ctx.Request.URL.Path}).Inc()
		// Histogram类型 meric的记录方式
		WebRequestDuration.With(prometheus.Labels{"method": ctx.Request.Method, "endpoint": ctx.Request.URL.Path}).Observe(duration.Seconds())
	}
}

//模拟查询
func Check(w http.ResponseWriter, r *http.Request) {
	//模拟业务查询耗时0~1s
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	w.Write([]byte(`{"alive": true}`))
}
