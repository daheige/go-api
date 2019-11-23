package main

import (
	"log"
	"sync"
	"time"

	httpRequest "github.com/daheige/thinkgo/gresty"
	"github.com/pkg/profile"
)

//查看fd情况 $ lsof -p pid -i | wc -l
func main() {
	defer profile.Start().Stop()

	s := &httpRequest.Service{
		BaseUri: "",
		Timeout: 2 * time.Second,
	}

	opt := &httpRequest.ReqOpt{
		Data: map[string]interface{}{
			"id": "1234",
		},
	}

	res := s.Do("post", "http://localhost:1338/v1/data", opt)
	log.Println("err: ", res.Err)
	log.Println("body:", res.Text())

	nums := 30000
	//每秒100个进行请求
	var wg sync.WaitGroup
	wg.Add(nums)
	for i := 0; i < nums; i++ {
		time.Sleep(10 * time.Millisecond)
		go func() {
			defer wg.Done()

			res := s.Do("post", "http://localhost:1338/v1/data", opt)
			log.Println("err: ", res.Err)
			log.Println("body:", res.Text())
		}()
	}

	wg.Wait()
	log.Println("ok")
}

/**
测试结果
2019/08/31 19:42:38 ok
2019/08/31 19:42:38 profile: cpu profiling disabled, /tmp/profile781753137/cpu.pprof

测试完毕后
$ lsof -p 9264 | wc -l
12

查看pprof
$ go tool pprof -http=:8080 /tmp/profile781753137/cpu.pprof
火焰图:  访问http://localhost:8080/ui/flamegraph
*/
