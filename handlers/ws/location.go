package ws

import (
	"encoding/json"
	"oui/models/ride"
)

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type LocationInfo struct {
	Name     string   `json:"name"`
	Location Location `json:"location"`
}

func NilLocation() (l Location) { return }

func OnCurrentLocation(c *Client, data interface{}) {
	jsonString, _ := json.Marshal(data)
	c.Location = Location{}
	json.Unmarshal(jsonString, &c.Location)
}

func OnAskTaxiLocation(c *Client, data interface{}) {
	taxiLocation := make([]LocationInfo, 0)
	for _, client := range TaxiUsers {
		if client != c && client.Location != NilLocation() {
			taxiLocation = append(taxiLocation, LocationInfo{client.User.Firstname + " " + client.User.Lastname, client.Location})
		}
	}
	c.Send("taxiLocation", taxiLocation)
}

func OnStopLocation(c *Client, data interface{}) {
	TaxiUsers[c.User.ID].Location = NilLocation()
}

func OnNewRide(c *Client, data interface{}) {
	jsonString, _ := json.Marshal(data)
	var r ride.Ride
	json.Unmarshal(jsonString, &r)
	r.Operator = c.User.ID
	ok := r.UpsertPgSQL()
	c.Send("newRide", map[string]interface{}{"success": ok})
}
