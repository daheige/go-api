# gin 1.4.0 大量请求panic问题 
    render/json.go 58行
    
    2019/11/23 16:38:46 http: panic serving 127.0.0.1:47148: write tcp 127.0.0.1:1338->127.0.0.1:47148: write: broken pipe
    goroutine 2695 [running]:
    net/http.(*conn).serve.func1(0xc00245b0e0)
    	/usr/local/go/src/net/http/server.go:1767 +0x139
    panic(0xbcfbc0, 0xc000d5cb40)
    	/usr/local/go/src/runtime/panic.go:679 +0x1b2
    github.com/gin-gonic/gin/render.JSON.Render(...)
    	/mygo/pkg/mod/github.com/gin-gonic/gin@v1.4.0/render/json.go:58
    github.com/gin-gonic/gin.(*Context).Render(0xc000dbd760, 0x1f4, 0xdd9740, 0xc0001cc1e0)
    	/mygo/pkg/mod/github.com/gin-gonic/gin@v1.4.0/context.go:768 +0x13f
    github.com/gin-gonic/gin.(*Context).JSON(...)
    	/mygo/pkg/mod/github.com/gin-gonic/gin@v1.4.0/context.go:811
    github.com/gin-gonic/gin.(*Context).AbortWithStatusJSON(...)
    	/mygo/pkg/mod/github.com/gin-gonic/gin@v1.4.0/context.go:155
    go-api/app/middleware.(*LogWare).Recover.func1.1(0xc000dbd760)
    	/web/go/go-api/app/middleware/LogWare.go:60 +0x299
    panic(0xbcfbc0, 0xc000d5cb40)
    	/usr/local/go/src/runtime/panic.go:679 +0x1b2
    github.com/gin-gonic/gin/render.JSON.Render(...)
    	/mygo/pkg/mod/github.com/gin-gonic/gin@v1.4.0/render/json.go:58
    github.com/gin-gonic/gin.(*Context).Render(0xc000dbd760, 0xc8, 0xdd9740, 0xc0024fbb80)
    	/mygo/pkg/mod/github.com/gin-gonic/gin@v1.4.0/context.go:768 +0x13f
    github.com/gin-gonic/gin.(*Context).JSON(...)
    	/mygo/pkg/mod/github.com/gin-gonic/gin@v1.4.0/context.go:811
    go-api/app/controller.(*BaseController).ajaxReturn(0x12d7628, 0xc000dbd760, 0x0, 0xca344a, 0x2, 0xb61ca0, 0xc001900d20)
    	/web/go/go-api/app/controller/BaseController.go:31 +0x2a5
    go-api/app/controller.(*BaseController).Success(...)
    	/web/go/go-api/app/controller/BaseController.go:44
    go-api/app/controller.(*IndexController).Hello(0x12d7628, 0xc000dbd760)
    	/web/go/go-api/app/controller/IndexController.go:32 +0x3af
    github.com/gin-gonic/gin.(*Context).Next(0xc000dbd760)
    	/mygo/pkg/mod/github.com/gin-gonic/gin@v1.4.0/context.go:124 +0x3b
    go-api/app/helper.Monitor.func1(0xc000dbd760)
    	/web/go/go-api/app/helper/utils.go:36 +0x51
    github.com/gin-gonic/gin.(*Context).Next(0xc000dbd760)
    	/mygo/pkg/mod/github.com/gin-gonic/gin@v1.4.0/context.go:124 +0x3b
    go-api/app/middleware.TimeoutHandler.func1(0xc000dbd760)
    	/web/go/go-api/app/middleware/ReqWare.go:48 +0x19a
    github.com/gin-gonic/gin.(*Context).Next(0xc000dbd760)
    	/mygo/pkg/mod/github.com/gin-gonic/gin@v1.4.0/context.go:124 +0x3b
    go-api/app/middleware.(*LogWare).Recover.func1(0xc000dbd760)
    	/web/go/go-api/app/middleware/LogWare.go:69 +0x5b
    github.com/gin-gonic/gin.(*Context).Next(0xc000dbd760)
    	/mygo/pkg/mod/github.com/gin-gonic/gin@v1.4.0/context.go:124 +0x3b
    go-api/app/middleware.(*LogWare).Access.func1(0xc000dbd760)
    	/web/go/go-api/app/middleware/LogWare.go:32 +0x150
    github.com/gin-gonic/gin.(*Context).Next(0xc000dbd760)
    	/mygo/pkg/mod/github.com/gin-gonic/gin@v1.4.0/context.go:124 +0x3b
    github.com/gin-gonic/gin.(*Engine).handleHTTPRequest(0xc000294280, 0xc000dbd760)
    	/mygo/pkg/mod/github.com/gin-gonic/gin@v1.4.0/gin.go:389 +0x5b2
    github.com/gin-gonic/gin.(*Engine).ServeHTTP(0xc000294280, 0xdddfc0, 0xc001f3ca80, 0xc001f99f00)
    	/mygo/pkg/mod/github.com/gin-gonic/gin@v1.4.0/gin.go:351 +0x134
    net/http.serverHandler.ServeHTTP(0xc0004420e0, 0xdddfc0, 0xc001f3ca80, 0xc001f99f00)
    	/usr/local/go/src/net/http/server.go:2802 +0xa4
    net/http.(*conn).serve(0xc00245b0e0, 0xde02c0, 0xc001968c80)
    	/usr/local/go/src/net/http/server.go:1890 +0x875
    created by net/http.(*Server).Serve
    	/usr/local/go/src/net/http/server.go:2927 +0x38e
    
    这里当写入流无法写入的时候，就不能再写入，当然这里抛出panic问题，导致其他接口响应时间变慢
    严重的时候会导致cpu暴涨，服务直接挂掉。