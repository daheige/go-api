package main

import (
	"context"
	"flag"
	"fmt"
	"go-api/app/config"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/profile"

	"github.com/daheige/thinkgo/monitor"

	"github.com/daheige/thinkgo/logger"

	"github.com/prometheus/client_golang/prometheus"

	"go-api/app/routes"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var port int
var log_dir string
var config_dir string
var wait time.Duration //平滑重启的等待时间1s or 1m

func init() {
	flag.IntVar(&port, "port", 1338, "app listen port")
	flag.StringVar(&log_dir, "log_dir", "./logs", "log dir")
	flag.StringVar(&config_dir, "config_dir", "./", "config dir")
	flag.DurationVar(&wait, "graceful-timeout", 3*time.Second, "the server gracefully reload. eg: 15s or 1m")
	flag.Parse()

	//日志文件设置
	logger.SetLogDir(log_dir)
	logger.SetLogFile("go-api.log")
	logger.MaxSize(500)
	logger.InitLogger()

	//初始化配置文件
	config.InitConf(config_dir)
	config.InitRedis()

	//注册监控指标
	prometheus.MustRegister(monitor.WebRequestTotal)
	prometheus.MustRegister(monitor.WebRequestDuration)
	prometheus.MustRegister(monitor.CpuTemp)
	prometheus.MustRegister(monitor.HdFailures)

	//性能监控的端口port+1000,只能在内网访问
	go func() {
		defer logger.Recover()

		pprof_port := port + 1000
		log.Println("server pprof run on: ", pprof_port)

		httpMux := http.NewServeMux() //创建一个http ServeMux实例
		httpMux.HandleFunc("/debug/pprof/", pprof.Index)
		httpMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		httpMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		httpMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		httpMux.HandleFunc("/debug/pprof/trace", pprof.Trace)
		httpMux.HandleFunc("/check", HealthCheckHandler)

		//metrics监控
		httpMux.Handle("/metrics", promhttp.Handler())

		if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", pprof_port), httpMux); err != nil {
			log.Println(err)
		}
	}()

	//gin mode设置
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
	//退出时候会自动采集profile性能指标

	//defer profile.Start(profile.MemProfile, profile.BlockProfile, profile.MutexProfile, profile.ThreadcreationProfile).Stop()

	defer profile.Start().Stop()

	router := gin.New()

	//加载路由文件中的路由
	routes.WebRoute(router)

	//服务server设置
	// IdleTimeout is the maximum amount of time to wait for the
	// next request when keep-alives are enabled. If IdleTimeout
	// is zero, the value of ReadTimeout is used. If both are
	// zero, ReadHeaderTimeout is used.

	// ReadHeaderTimeout is the amount of time allowed to read
	// request headers. The connection's read deadline is reset
	// after reading the headers and the Handler can decide what
	// is considered too slow for the body.

	//对于idleTimeout一般不建议设置，如果不设置默认采用ReadTimeout
	//对于ReadHeaderTimeout一般不建议设置，默认采用ReadTimeout
	//详细分析: https://blog.csdn.net/busai2/article/details/82634049
	server := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
		//IdleTimeout:       10 * time.Second, //tcp idle time
		//ReadHeaderTimeout: 10 * time.Second, //read header timeout
		ReadTimeout:  5 * time.Second,  //read request timeout
		WriteTimeout: 10 * time.Second, //write timeout
	}

	//在独立携程中运行
	log.Println("server run on: ", port)
	log.Println("server pid: ", os.Getppid())

	go func() {
		defer logger.Recover()

		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	//平滑重启
	ch := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// recivie signal to exit main goroutine
	// window signal
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2, os.Interrupt, syscall.SIGHUP)

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
	go server.Shutdown(ctx) //在独立的携程中关闭服务器
	<-ctx.Done()

	log.Println("shutting down")
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	w.Write([]byte(`{"alive": true}`))
}
