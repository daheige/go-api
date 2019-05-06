package logic

import "log"

type HomeLogic struct {
	BaseLogic
}

func (h *HomeLogic) GetData() []string {
	log.Println(h.Ctx.GetString("current_uid"))

	return []string{
		"golang",
		"php",
	}
}
