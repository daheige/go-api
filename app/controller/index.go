package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type IndexController struct {
	BaseController
}

// Hello 压测
func (ctrl *IndexController) Hello(ctx *gin.Context) {
	// log.Println(111)
	userInfo := getUserInfo()

	// 模拟赋值操作
	info := make([]UserInfo, 0, len(userInfo))
	for k, _ := range userInfo {
		info = append(info, UserInfo{
			Id:      userInfo[k].Id,
			Name:    userInfo[k].Name,
			Age:     userInfo[k].Age,
			Content: userInfo[k].Content,
		})
	}

	// 对map进行主动gc
	userInfo = nil

	ctrl.Success(ctx, "ok", info)
}

// UserInfo user info.
type UserInfo struct {
	Id      int64
	Name    string
	Age     int
	Content string
}

func getUserInfo() map[string]UserInfo {
	user := make(map[string]UserInfo, 500)
	str := `What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
	What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
	What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
	What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
`

	for i := 0; i < 100; i++ {
		nick := "hello_" + strconv.Itoa(i)
		user[nick] = UserInfo{
			Id:      int64(i),
			Name:    nick,
			Age:     i + 10,
			Content: str,
		}
	}

	return user
}
