package controller

import (
	"go-api/app/logic"

	"github.com/gin-gonic/gin"
)

type HomeController struct {
	BaseController
}

// action
func (ctrl *HomeController) Index(ctx *gin.Context) {
	ctx.JSON(HTTP_SUCCESS_CODE, gin.H{
		"code":    200,
		"message": "ok",
		"data":    EmptyArray{},
	})
}

func (ctrl *HomeController) Test(ctx *gin.Context) {
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

func (ctrl *HomeController) GetData(ctx *gin.Context) {
	homeLogic := logic.HomeLogic{}
	homeLogic.SetCtx(ctx)

	data := homeLogic.GetData()

	ctx.JSON(HTTP_SUCCESS_CODE, gin.H{
		"code":    0,
		"message": "ok",
		"data":    data,
	})
}
