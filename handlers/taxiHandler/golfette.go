package taxiHandler

import (
	"oui/models"
	"oui/models/ride"

	"github.com/kataras/iris/v12"
)

func Golfette(c iris.Context, route models.Route) {

	if c.Method() != "GET" || route.Secondary == "" || route.Tertiary != "" {
		c.StopWithStatus(iris.StatusNotFound)
		return
	}
	switch route.Secondary {
	case "status":
		c.StopWithJSON(iris.StatusOK, iris.Map{"status": ride.GolfetteStatus})
	case "update":
		var q GolfetteQuery
		err := c.ReadQuery(&q)
		if err != nil || q.Token != "OpxY96T1ElSmAjng4xkocuBPQ8Kd" {
			c.StopWithStatus(iris.StatusBadRequest)
			return
		}
		ride.GolfetteStatus = q.Status
		c.StopWithJSON(iris.StatusOK, iris.Map{"status": ride.GolfetteStatus})
	default:
		c.StopWithStatus(iris.StatusNotFound)
	}

}
