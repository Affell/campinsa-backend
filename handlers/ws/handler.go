package ws

import (
	"net/http"
	"oui/models/ride"
	"oui/models/user"

	"github.com/gorilla/websocket"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
)

var TaxiUsers map[int64]*Client = make(map[int64]*Client)

type Handler func(*Client, interface{})

type Event string

type Router struct {
	rules map[Event]Handler
}

func NewRouter() *Router {
	return &Router{
		rules: make(map[Event]Handler),
	}
}

func (rt *Router) ServeHTTP(c iris.Context) {
	r := c.Request()
	w := c.ResponseWriter()
	var u user.User
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			token := c.Params().GetDefault("token", "").(string)
			if token == "" {
				golog.Debug("No token in header")
				return false
			}
			var err error
			if u, err = user.GetUserByTaxiToken(token); err != nil {
				golog.Debug("Invalid token in websocket : %v", err)
				return false
			}

			if old, ok := TaxiUsers[u.ID]; ok {
				old.socket.Close()
				old.Location = NilLocation()
			}

			return true
		},
	}

	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	client := NewClient(socket, rt.FindHandler, u)
	TaxiUsers[u.ID] = client

	if ride, ok := ride.Riders[u.ID]; ok {
		client.Send("rideAnswer", map[string]interface{}{"success": true, "ride": ride})
	}

	client.Read()
}

func (rt *Router) FindHandler(event Event) (Handler, bool) {
	handler, found := rt.rules[event]
	return handler, found
}

func (rt *Router) On(event Event, handler Handler) {
	rt.rules[event] = handler
}

func Broadcast(name string, data interface{}) {
	for _, c := range TaxiUsers {
		c.Send(name, data)
	}
}
