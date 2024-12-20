package ws

import (
	"encoding/json"
	"os"
	"oui/models/planning"
	"oui/models/ride"
	"oui/models/user"

	"github.com/kataras/golog"
)

func OnMode(c *Client, data interface{}) {
	jsonString, _ := json.Marshal(data)
	jsonData := struct {
		Type  string `json:"mode"`
		Token string `json:"token"`
	}{}
	json.Unmarshal(jsonString, &jsonData)

	if jsonData.Type == "external" {
		c.Send("authenticated", map[string]string{"tel": os.Getenv("STANDARD")})
		return
	}

	if jsonData.Type != "internal" || jsonData.Token == "" {
		c.socket.Close()
		return
	}

	u := user.User{ID: -1}
	var err error
	if u, err = user.GetUserByTaxiToken(jsonData.Token); err != nil {
		c.socket.Close()
		return
	}

	if old, ok := TaxiUsers[u.ID]; ok && old.socket != c.socket {
		old.Send("close", nil)
		old.socket.Close()
		old.Location = NilLocation()
	}

	c.User = u
	TaxiUsers[u.ID] = c
	if ride, ok := ride.Riders[u.ID]; ok {
		c.Send("rideAnswer", map[string]interface{}{"success": true, "ride": ride})
	}

	golog.Infof("CariTaxi authentification by : %v %v", u.Firstname, u.Lastname)
	c.Send("authenticated", nil)
}

func OnPlanning(c *Client, data interface{}) {
	if p, ok := planning.GlobalPlanning[c.User.ID]; ok {
		c.Send("planning", map[string]interface{}{"success": true, "planning": p})
	} else {
		c.Send("planning", map[string]interface{}{"success": false})
	}
}
