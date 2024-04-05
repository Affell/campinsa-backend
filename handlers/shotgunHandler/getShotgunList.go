package shotgunHandler

import (
	"oui/models"
	"oui/models/shotgun"
	"time"

	"github.com/kataras/iris/v12"
)

func List(c iris.Context, route models.Route) {

	data := struct {
		Shotguns []shotgun.Shotgun `json:"shotguns"`
	}{}

	shotguns, err := shotgun.GetAllShotguns()
	if err != nil {
		c.StopWithStatus(iris.StatusInternalServerError)
		return
	} else if len(shotguns) == 0 {
		c.StopWithStatus(iris.StatusNoContent)
		return
	}

	for _, s := range shotguns {
		if time.Now().UnixMilli() < s.UnlockTime {
			s.FormLink = ""
		}
		data.Shotguns = append(data.Shotguns, s)
	}

	c.StopWithJSON(iris.StatusOK, data)
}
