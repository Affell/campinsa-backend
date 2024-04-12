package ws

import (
	"encoding/json"
	"math"
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

func Distance(first, second Location) float64 {
	radlat1 := float64(math.Pi * first.Latitude / 180)
	radlat2 := float64(math.Pi * second.Latitude / 180)

	theta := float64(first.Longitude - second.Longitude)
	radtheta := float64(math.Pi * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)
	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / math.Pi
	dist = dist * 60 * 1.1515 * 1.609344

	return dist
}

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
