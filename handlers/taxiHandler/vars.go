package taxiHandler

import "oui/models"

const Service = "taxi"

var (
	Handlers models.HandlerMap = models.HandlerMap{
		"login":    Login,
		"download": Download,
		"update":   OnUpdate,
	}
)

type (
	LoginForm struct {
		Token string `form:"token" json:"token"`
	}
	GolfetteQuery struct {
		Token  string `form:"token" json:"token"`
		Status bool   `form:"status" json:"status"`
	}
)
