package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/daheige/go-api/config"
	// "github.com/pkg/profile"

	"github.com/daheige/thinkgo/gpprof"
	"github.com/daheige/thinkgo/logger"
	"github.com/daheige/thinkgo/monitor"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "go.uber.org/automaxprocs"

	"github.com/daheige/go-api/app/routes"
)

var (
	port      int
	logDir    string
	configDir string
	wait      time.Duration // 平滑重启的等待时间1s or 1m
)

func init() {
	flag.IntVar(&port, "port", 1338, "app listen port")
	flag.StringVar(&logDir, "log_dir", "./logs", "log dir")
	flag.StringVar(&configDir, "config_dir", "./", "config dir")
	flag.DurationVar(&wait, "graceful-timeout", 3*time.Second, "the server gracefully reload. eg: 15s or 1m")
	flag.Parse()

	// 日志文件设置
	logger.SetLogDir(logDir)
	logger.SetLogFile("go-api.log")
	logger.MaxSize(500)

	// 由于app/extensions/logger基于thinkgo/logger又包装了一层，所以这里是1
	logger.InitLogger(3)

	// 初始化配置文件
	config.InitConf(configDir)
	config.InitRedis()

	// 添加prometheus性能监控指标
	prometheus.MustRegister(monitor.WebRequestTotal)
	prometheus.MustRegister(monitor.WebRequestDuration)

	prometheus.MustRegister(monitor.CpuTemp)
	prometheus.MustRegister(monitor.HdFailures)

	// 性能监控的端口port+1000,只能在内网访问
	httpMux := gpprof.New()

	// 添加prometheus metrics处理器
	httpMux.Handle("/metrics", promhttp.Handler())
	gpprof.Run(httpMux, port+1000)

	// gin mode设置
	switch config.AppEnv {
	case "local", "dev":
		gin.SetMode(gin.DebugMode)
	case "testing":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}
}

func main() {
	// 开发环境可以打开这个profile性能分析
	// 退出时候会自动采集profile性能指标

	// defer profile.Start(profile.MemProfile, profile.BlockProfile,
	// 	profile.MutexProfile, profile.ThreadcreationProfile).Stop()

	// defer profile.Start().Stop()

	router := gin.New()

	// 加载路由文件中的路由
	routes.WebRoute(router)

	// 服务server设置
	// IdleTimeout is the maximum amount of time to wait for the
	// next request when keep-alives are enabled. If IdleTimeout
	// is zero, the value of ReadTimeout is used. If both are
	// zero, ReadHeaderTimeout is used.

	// ReadHeaderTimeout is the amount of time allowed to read
	// request headers. The connection's read deadline is reset
	// after reading the headers and the Handler can decide what
	// is considered too slow for the body.

	// 对于idleTimeout一般不建议设置，如果不设置默认采用ReadTimeout
	// 对于ReadHeaderTimeout一般不建议设置，默认采用ReadTimeout
	// 详细分析: https://blog.csdn.net/busai2/article/details/82634049
	server := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 2 << 20, // header max 2MB
	}

	// 在独立携程中运行
	log.Println("server run on: ", port)
	log.Println("server pid: ", os.Getppid())

	go func() {
		defer logger.Recover()

		if err := server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				logger.Info("server close error", map[string]interface{}{
					"trace_error": err.Error(),
				})

				log.Println(err)
				return
			}

			log.Println("server will exit...")
		}
	}()

	// 平滑重启
	ch := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// recivie signal to exit main goroutine
	// window signal
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)

	// linux signal if you use linux on production,please use this code.
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2, os.Interrupt, syscall.SIGHUP)

	// Block until we receive our signal.
	sig := <-ch

	log.Println("exit signal: ", sig.String())
	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// if your application should wait for other services
	// to finalize based on context cancellation.
	go server.Shutdown(ctx)
	<-ctx.Done()

	log.Println("server shutting down")
}
