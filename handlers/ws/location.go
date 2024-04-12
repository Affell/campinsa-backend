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

func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

func Distance(first, second Location) float64 {
	radlat1 := float64(math.Pi * first.Latitude / 180)
	radlng1 := float64(math.Pi * first.Longitude / 180)
	radlat2 := float64(math.Pi * second.Latitude / 180)
	radlng2 := float64(math.Pi * second.Longitude / 180)

	r := 6378.1
	h := hsin(radlat2-radlat1) + math.Cos(radlat1)*math.Cos(radlat2)*hsin(radlng2-radlng1)
	return 2 * r * math.Asin(math.Sqrt(h))
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
