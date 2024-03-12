package ws

import (
	"encoding/json"
	"oui/models/ride"
	"strconv"
)

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
		Broadcast("newRide", m)
	} else {
		c.Send("newRide", m)
	}
}

func RetrieveRides(c *Client, data interface{}) {
	rides := ride.GetAllRides()
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
		Broadcast("updateRide", map[string]interface{}{"ride": r})
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
		Broadcast("updateRide", map[string]interface{}{"ride": r})
	}
}
