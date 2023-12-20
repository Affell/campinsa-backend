package shotgunHandler

import (
	"oui/models"
	"oui/models/shotgun"
	"oui/models/user"
	"time"

	"github.com/kataras/iris/v12"
)

func List(c iris.Context, route models.Route) {

	var token user.UserToken
	if t := c.GetID(); t != nil {
		token = t.(user.UserToken)
	}

	data := struct {
		Running []shotgun.Shotgun `json:"running"`
		Ended   []shotgun.Shotgun `json:"ended"`
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
		if !user.HasPermission(token.ID, "shotgun.detail") {
			s.CreatedTime = -1
			if time.Now().UnixMilli() < s.UnlockTime {
				s.FormLink = ""
			}
		}
		if s.Ended {
			data.Ended = append(data.Ended, s)
		} else {
			data.Running = append(data.Running, s)
		}
	}

	c.StopWithJSON(iris.StatusOK, data)
}
