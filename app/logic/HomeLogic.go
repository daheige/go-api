package logic

import (
	"go-api/app/model"
	"log"

	"github.com/daheige/thinkgo/mysql"
)

type HomeLogic struct {
	BaseLogic
}

func (h *HomeLogic) GetData(id string) []string {
	log.Println(h.Ctx.GetString("current_uid"))

	db, err := mysql.GetDbObj("default")
	if err != nil {
		log.Println("db connection error: ", err)
		return nil
	}

	user := &model.User{}
	db.Where("name = ?", "hello").First(user)
	log.Println("user: ", user)

	if id == "1234" {
		return []string{
			"js",
			"php",
			user.Name,
		}
	}

	return []string{
		"golang",
		"php",
	}
}
