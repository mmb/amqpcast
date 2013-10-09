package amqpcast

import (
	"code.google.com/p/go.net/websocket"
)

type Connection struct {
	Ws       *websocket.Conn
	Outbound chan string
}

func (c *Connection) write() {
	for message := range c.Outbound {
		err := websocket.Message.Send(c.Ws, message)
		if err != nil {
			break
		}
	}
}

func (c *Connection) read() {
	for {
		var message string
		err := websocket.Message.Receive(c.Ws, &message)
		if err != nil {
			break
		}
	}
}
