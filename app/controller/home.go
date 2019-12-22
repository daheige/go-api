package controller

import (
	"go-api/app/extensions/logger"
	"go-api/app/logic"
	"log"
	"time"

	"go-api/app/config"

	"github.com/gin-gonic/gin"
)

type HomeController struct {
	BaseController
}

// action
func (ctrl *HomeController) Index(ctx *gin.Context) {

	logger.Info(ctx.Request.Context(), "1234fe", nil)

	ctx.JSON(HTTP_SUCCESS_CODE, gin.H{
		"code":    200,
		"message": "ok",
	})
}

func (ctrl *HomeController) Test(ctx *gin.Context) {
	panic(11)

	ctx.JSON(HTTP_SUCCESS_CODE, gin.H{
		"code":    0,
		"message": "ok",
		"data":    "this is test",
	})
}

func (ctrl *HomeController) Info(ctx *gin.Context) {
	id := ctx.Param("id")
	ctx.JSON(HTTP_SUCCESS_CODE, gin.H{
		"code":    0,
		"message": "ok",
		"data":    "current id: " + id,
	})
}

func (ctrl *HomeController) GetInfo(ctx *gin.Context) {
	ctx.JSON(HTTP_SUCCESS_CODE, gin.H{
		"code":    0,
		"message": "ok",
		"data": map[string]interface{}{
			"id":   1,
			"name": "heige",
			"sex":  1,
		},
	})
}

// GetData 模拟数据库查询数据
func (ctrl *HomeController) GetData(ctx *gin.Context) {
	homeLogic := &logic.HomeLogic{}
	homeLogic.SetCtx(ctx)
	name := ctx.DefaultQuery("name", "hello")

	if name == "" {
		ctx.JSON(HTTP_SUCCESS_CODE, gin.H{
			"code":    500,
			"message": "name is empty",
		})

		return
	}

	data, err := homeLogic.GetData(name)
	if err != nil {
		ctx.JSON(HTTP_SUCCESS_CODE, gin.H{
			"code":    500,
			"message": "get user fail: " + err.Error(),
		})

		return
	}

	ctx.JSON(HTTP_SUCCESS_CODE, gin.H{
		"code":    0,
		"message": "ok",
		"data":    data,
	})
}

// PostData 模拟post数据
func (ctrl *HomeController) PostData(ctx *gin.Context) {
	homeLogic := &logic.HomeLogic{}
	homeLogic.SetCtx(ctx)
	name := ctx.DefaultPostForm("name", "hello")

	data, err := homeLogic.GetData(name)
	log.Println("err: ", err)

	ctx.JSON(HTTP_SUCCESS_CODE, gin.H{
		"code":    0,
		"message": "ok",
		"data":    data,
	})
}

// SetData 模拟redis数据设置
func (ctrl *HomeController) SetData(ctx *gin.Context) {
	redisObj, err := config.GetRedisObj("default")
	if err != nil {
		log.Println(err)
		ctx.JSON(200, gin.H{
			"code":    500,
			"message": "redis connection error",
		})

		return
	}

	//用完就需要释放连接，防止过多的连接导致redis连接过多而陷入长久等待，从而redis崩溃
	defer redisObj.Close()

	_, err = redisObj.Do("set", "myname", "daheige")
	if err != nil {
		logger.Error(ctx.Request.Context(), "set redis error", map[string]interface{}{
			"trace_error": err.Error(),
		})

		ctx.JSON(200, gin.H{
			"code":    500,
			"message": "set data error",
		})

		return
	}

	ctx.JSON(200, gin.H{
		"code":    0,
		"message": "set data success",
	})
}

// LongAsync When starting new Goroutines inside a middleware or handler,
// you SHOULD NOT use the original context inside it,
// you have to use a read-only copy.
func (ctrl *HomeController) LongAsync(ctx *gin.Context) {
	// create copy to be used inside the goroutine
	cCp := ctx.Copy()
	go func() {
		// simulate a long task with time.Sleep(). 3 seconds
		time.Sleep(3 * time.Second)

		// note that you are using the copied context "cCp", IMPORTANT
		log.Println("Done! in path " + cCp.Request.URL.Path)
	}()

	ctx.JSON(200, gin.H{
		"code":    0,
		"message": "ok",
	})
}

type Person struct {
	Id      int64  `form:"id" binding:"required,min=1"`
	Name    string `form:"name" binding:"omitempty"` //可选参数
	Address string `form:"address" binding:"required"`
}

// ValidData 测试gin(1.5.0+) binding功能
// http://localhost:1338/v1/person-info?id=0&address=fefefe 参数错误
// http://localhost:1338/v1/person-info?id=12&address=fefefe 参数符合预期
func (ctrl *HomeController) ValidData(ctx *gin.Context) {
	p := &Person{}
	if err := ctx.ShouldBind(p); err != nil {
		log.Println("error: ", err)

		ctrl.Error(ctx, 500, "param error")
		return
	}

	log.Println("id: ", p.Id)
	log.Println("name: ", p.Name, "address: ", p.Address)

	ctrl.Success(ctx, "ok", p)
}
