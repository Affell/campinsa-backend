package shotgunHandler

import (
	"oui/models"
	"oui/models/shotgun"
	"time"

	"github.com/kataras/iris/v12"
)

func List(c iris.Context, route models.Route) {

	data := struct {
		Thursday map[string]interface{}   `json:"thursday"`
		Friday   map[string]interface{}   `json:"friday"`
		Shotguns []map[string]interface{} `json:"shotguns"`
	}{}

	shotguns, err := shotgun.GetAllShotguns()
	if err != nil {
		c.StopWithStatus(iris.StatusInternalServerError)
		return
	}
	now := time.Now()
	for _, s := range shotguns {
		m := s.ToWebDetails()
		if s.Id == "1" {
			data.Friday = m
		} else if s.Id == "2" {
			data.Thursday = m
		} else if now.Weekday() == time.Unix(0, s.UnlockTime*int64(time.Second)).Weekday() {
			data.Shotguns = append(data.Shotguns, m)
		}
	}

	c.StopWithJSON(iris.StatusOK, data)
}
