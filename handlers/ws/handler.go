package ws

import (
	"net/http"
	"oui/models/user"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
)

var TaxiUsers sync.Map

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
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			token := c.Params().GetDefault("token", "").(string)
			if token == "" {
				golog.Debug("No token in header")
				return false
			}
			var u user.User
			var err error
			if u, err = user.GetUserByTaxiToken(token); err != nil {
				golog.Debug("Invalid token in websocket : %v", err)
				return false
			}

			TaxiUsers.Store(c.GetID(), u)
			return true
		},
	}

	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, ok := TaxiUsers.LoadAndDelete(c.GetID())
	if !ok {
		golog.Error("Unable to load TaxiUser infos")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	client := NewClient(socket, rt.FindHandler)
	u := data.(user.User)
	client.User = u
	TaxiUsers.Store(c.GetID(), client)

	client.Read()
}

func (rt *Router) FindHandler(event Event) (Handler, bool) {
	handler, found := rt.rules[event]
	return handler, found
}

func (rt *Router) On(event Event, handler Handler) {
	rt.rules[event] = handler
}
