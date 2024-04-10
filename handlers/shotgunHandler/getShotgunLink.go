package shotgunHandler

import (
	"oui/models"
	"oui/models/shotgun"
	"time"

	"github.com/kataras/iris/v12"
)

func GetLink(c iris.Context, route models.Route) {

	if c.Method() != "GET" || route.Secondary == "" || route.Tertiary != "" || len(route.Tail) != 0 {
		c.StopWithStatus(iris.StatusNotFound)
		return
	}

	s, err := shotgun.GetShotgun(route.Secondary)
	if err != nil {
		c.StopWithStatus(iris.StatusNotFound)
		return
	}

	if time.Now().Unix() < s.UnlockTime {
		c.StopWithStatus(iris.StatusForbidden)
		return
	}

	c.Redirect(s.FormLink, iris.StatusFound)
}
