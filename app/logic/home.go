package logic

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/daheige/go-api/app/model"

	"github.com/daheige/thinkgo/mysql"
	"github.com/jinzhu/gorm"
)

// HomeLogic home logic.
type HomeLogic struct {
	BaseLogic
}

// GetData 模拟数据库查询
func (h *HomeLogic) GetData(name string) (map[string]interface{}, error) {
	db, err := mysql.GetDbObj("default")
	if err != nil {
		// log.Println("db connection error: ", err)
		return nil, errors.New("db connection error")
	}

	user := &model.User{}
	err = db.Where("name = ?", name).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	// log.Println("user: ", user)

	return map[string]interface{}{
		"user": user,
	}, nil
}

// AsyncDoTaskByCtx 通过ctx控制执行任务
func (h *HomeLogic) AsyncDoTaskByCtx(ctx context.Context, id int) {
	done := make(chan struct{}, 1)
	go func() {
		name := ctx.Value("name")
		log.Printf("name:%s\n", name)
		for i := 0; i < id; i++ {
			log.Println("current index: ", i)
		}

		time.Sleep(3 * time.Second) // 模拟超时
		log.Println(1111111)
		close(done)
	}()

	select {
	case <-ctx.Done(): // ctx超时时间
		log.Println("ctx timeout,error: ", ctx.Err())
	case <-done: // 业务内部指定的一个超时时间
		log.Println("success")
	}
}
