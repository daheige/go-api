package logic

import "github.com/gin-gonic/gin"

type BaseLogic struct {
	Ctx *gin.Context
}

func (b *BaseLogic) SetCtx(ctx *gin.Context) {
	b.Ctx = ctx
}
