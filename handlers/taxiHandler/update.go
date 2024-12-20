package taxiHandler

import (
	"oui/models"
	"oui/models/planning"

	"github.com/kataras/iris/v12"
)

func OnUpdate(c iris.Context, route models.Route) {
	if c.Method() != "GET" || route.Secondary != "" {
		c.StopWithStatus(iris.StatusNotFound)
		return
	}
	var body struct {
		Token string
	}
	err := c.ReadBody(&body)
	if err != nil || body.Token != planning.UpdateToken {
		c.StopWithStatus(iris.StatusUnauthorized)
		return
	}

	planning.UpdatePlanning()

	c.StopWithJSON(iris.StatusOK, iris.Map{"success": true})

}
