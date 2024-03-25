package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/kataras/iris/v12"
)

var TaxiUsers map[int64]*Client = make(map[int64]*Client)

type HandlerDesc struct {
	HandlerFunc  HandlerFunc
	AuthRequired bool
}
type HandlerFunc func(*Client, interface{})

type Event string

type Router struct {
	rules map[Event]HandlerDesc
}

func NewRouter() *Router {
	return &Router{
		rules: make(map[Event]HandlerDesc),
	}
}

func (rt *Router) ServeHTTP(c iris.Context) {
	r := c.Request()
	w := c.ResponseWriter()
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	client := NewClient(socket, rt.FindHandler)
	client.Send("mode", nil)
	client.Read()
}

func (rt *Router) FindHandler(event Event) (HandlerDesc, bool) {
	handler, found := rt.rules[event]
	return handler, found
}

func (rt *Router) On(event Event, handler HandlerFunc, authenticated bool) {
	rt.rules[event] = HandlerDesc{HandlerFunc: handler, AuthRequired: authenticated}
}

func Broadcast(name string, data interface{}, auth bool) {
	for _, c := range TaxiUsers {
		if !auth || c.User.ID != -1 {
			c.Send(name, data)
		}
	}
}
