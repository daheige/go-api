package logic

import (
	"context"
)

// BaseLogic
type BaseLogic struct {
	Ctx context.Context
}

// SetCtx set ctx
func (b *BaseLogic) SetCtx(ctx context.Context) {
	b.Ctx = ctx
}
