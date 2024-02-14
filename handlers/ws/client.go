package ws

import (
	"oui/models/user"

	"github.com/gorilla/websocket"
	"github.com/kataras/golog"
)

type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

type FindHandler func(Event) (Handler, bool)

type Client struct {
	socket      *websocket.Conn
	findHandler FindHandler
	User        user.User
	Location    Location
}

func NewClient(socket *websocket.Conn, findHandler FindHandler, user user.User) *Client {
	return &Client{
		socket:      socket,
		findHandler: findHandler,
		User:        user,
	}
}

func (c *Client) Send(name string, data interface{}) {
	msg := Message{name, data}
	err := c.socket.WriteJSON(msg)
	if err != nil {
		golog.Errorf("socket write error: %v\n", err)
	}
}

func (c *Client) Read() {
	var msg Message
	for {
		if err := c.socket.ReadJSON(&msg); err != nil {
			break
		}
		if handler, found := c.findHandler(Event(msg.Name)); found {
			handler(c, msg.Data)
		}
	}

	delete(TaxiUsers, c.User.ID)

	c.socket.Close()
}