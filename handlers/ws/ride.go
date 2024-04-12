package ws

import (
	"encoding/json"
	"oui/models/ride"
	"strconv"
	"time"

	"github.com/kataras/golog"
)

func BroadcastNewRide(r ride.Ride) {

	message := map[string]interface{}{
		"success": true,
		"ride":    r,
	}

	var targets []*Client
	loc := Location{
		Latitude:  r.Start.Latitude,
		Longitude: r.Start.Longitude,
	}
	m := make(map[int64]float64)
	for _, user := range TaxiUsers {
		if user.User.ID != 0 && user.Location != NilLocation() {
			m[user.User.ID] = Distance(loc, user.Location)
			if len(targets) == 0 {
				targets = append(targets, user)
			} else {
				for i, t := range targets {
					if m[user.User.ID] < m[t.User.ID] {
						targets = insert(targets, i, user)
						break
					}
				}
			}
		} else {
			golog.Debugf("%v %v", user.User.ID, user.Location)
		}
	}

	if len(targets) > 0 {
		targets[0].Send("newRide", message)
		if len(targets) > 1 {
			time.AfterFunc(time.Second*5, func() {
				update, err := ride.GetRideByID(r.ID)
				if err != nil {
					Broadcast("newRide", message, true)
					return
				}
				if update.Taxi == 0 {
					targets[1].Send("newRide", message)
				}
				if len(targets) > 2 {
					time.AfterFunc(time.Second*5, func() {
						update, err := ride.GetRideByID(r.ID)
						if err != nil || update.Taxi == 0 {
							for i := 2; i < len(targets); i++ {
								targets[i].Send("newRide", message)
							}
						}
					})
				}
			})
		}
	}
}

func insert(a []*Client, index int, value *Client) []*Client {
	if len(a) == index {
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...)
	a[index] = value
	return a
}

func OnNewRide(c *Client, data interface{}) {
	jsonString, _ := json.Marshal(data)
	var r ride.Ride
	json.Unmarshal(jsonString, &r)
	r.Operator = c.User.ID
	r.TranslateRideIds()
	ok := r.UpsertPgSQL()
	m := map[string]interface{}{"success": ok}
	if ok {
		m["ride"] = r.ToAppDetails()
		c.Send("newRide", m)
		if r.Date != 0 {
			time.AfterFunc(time.Until(time.Unix(r.Date/1000, 0)), func() {
				BroadcastNewRide(r)
			})
		} else {
			BroadcastNewRide(r)
		}
	} else {
		c.Send("newRide", m)
	}
}

func RetrieveRides(c *Client, data interface{}) {
	rides := ride.GetAllRides(false)
	c.Send("rides", map[string]interface{}{
		"rides": rides,
	})
}

func OnRideAnswer(c *Client, data interface{}) {
	jsonString, _ := json.Marshal(data)
	var d map[string]interface{}
	json.Unmarshal(jsonString, &d)
	if strId, ok := d["id"]; ok {
		id, err := strconv.ParseInt(strId.(string), 10, 64)
		if err != nil {
			c.Send("rideAnswer", map[string]interface{}{"success": false, "message": "Course indisponible"})
			return
		}
		r, err := ride.GetRideByID(id)
		if err != nil || r.Taxi != 0 {
			c.Send("rideAnswer", map[string]interface{}{"success": false, "message": "Course indisponible"})
			return
		}
		if _, ok := ride.Riders[c.User.ID]; ok {
			c.Send("rideAnswer", map[string]interface{}{"success": false, "message": "Vous avez déjà une course en cours"})
			return
		}
		r.Taxi = c.User.ID
		r.UpsertPgSQL()
		r.TranslateRideIds()
		ride.Riders[c.User.ID] = r
		c.Send("rideAnswer", map[string]interface{}{"success": true, "ride": r})
		Broadcast("updateRide", map[string]interface{}{"ride": r}, true)
	}
}

func OnRideCompleted(c *Client, data interface{}) {
	jsonString, _ := json.Marshal(data)
	var d map[string]interface{}
	json.Unmarshal(jsonString, &d)
	if strId, ok := d["id"]; ok {
		id, err := strconv.ParseInt(strId.(string), 10, 64)
		if err != nil {
			c.Send("rideCompleted", map[string]interface{}{"success": false, "message": "Course indisponible"})
			return
		}
		r, err := ride.GetRideByID(id)
		if err != nil || r.Taxi != c.User.ID {
			c.Send("rideCompleted", map[string]interface{}{"success": false, "message": "Vous n'êtes pas le conducteur assigné à cette course"})
		}
		r.Completed = true
		r.UpsertPgSQL()
		delete(ride.Riders, c.User.ID)
		c.Send("rideCompleted", map[string]interface{}{"success": true})
		Broadcast("updateRide", map[string]interface{}{"ride": r}, true)
	}
}
