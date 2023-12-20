package ws

import (
	"encoding/json"
)

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type LocationInfo struct {
	Name     string   `json:"name"`
	Location Location `json:"location"`
}

func OnCurrentLocation(c *Client, data interface{}) {
	jsonString, _ := json.Marshal(data)
	c.Location = Location{}
	json.Unmarshal(jsonString, &c.Location)
	taxiLocation := make([]LocationInfo, 0)
	TaxiUsers.Range(func(key, c any) bool {
		client := c.(*Client)
		if client != c {
			taxiLocation = append(taxiLocation, LocationInfo{client.User.Firstname + " " + client.User.Lastname, client.Location})
		}
		return true
	})
	c.Send("taxiLocation", taxiLocation)
}
