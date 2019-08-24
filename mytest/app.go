package main

import (
	"log"
	"sync"
	"time"

	"github.com/daheige/thinkgo/httpRequest"
)

//查看fd情况 $ lsof -p pid -i | wc -l
func main() {
	s := &httpRequest.Service{
		BaseUri: "",
		Timeout: 2,
	}

	opt := &httpRequest.ReqOpt{
		Params: map[string]interface{}{
			"id": "1234",
		},
	}

	nums := 30000
	//每秒100个进行请求
	var wg sync.WaitGroup
	wg.Add(nums)
	for i := 0; i < nums; i++ {
		time.Sleep(10 * time.Millisecond)
		go func() {
			defer wg.Done()

			res := s.Do("get", "http://localhost:1338/v1/data", opt)
			log.Println("err: ", res.Err)
			log.Println("body:", res.Text())
		}()
	}

	wg.Wait()
	log.Println("ok")
}
