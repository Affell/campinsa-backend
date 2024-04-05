package taxiHandler

import "oui/models"

const Service = "taxi"

var (
	Handlers models.HandlerMap = models.HandlerMap{
		"login":    Login,
		"download": Download,
	}
)

type (
	LoginForm struct {
		Token string `form:"token" json:"token"`
	}
)
