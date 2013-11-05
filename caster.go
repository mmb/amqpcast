package amqpcast

import (
	"log"
	"sync"
)

type Caster struct {
	sync.RWMutex
	Connections map[*Connection]bool
	Create      chan *Connection
	Destroy     chan *Connection
	Outbound    chan string
}

func NewCaster() *Caster {
	return &Caster{
		Connections: make(map[*Connection]bool),
		Create:      make(chan *Connection),
		Destroy:     make(chan *Connection),
		Outbound:    make(chan string, 256),
	}
}

func (cstr *Caster) Run() {
	for {
		select {
		case c := <-cstr.Create:
			log.Printf("new client")
			cstr.createConnection(c)
		case c := <-cstr.Destroy:
			log.Printf("client closed")
			cstr.destroyConnection(c)
		case m := <-cstr.Outbound:
			cstr.distributeMessage(m)
		}
	}
}

func (cstr *Caster) createConnection(c *Connection) {
	cstr.Lock()
	cstr.Connections[c] = true
	cstr.Unlock()
}

func (cstr *Caster) destroyConnection(c *Connection) {
	cstr.Lock()
	delete(cstr.Connections, c)
	cstr.Unlock()
	c.close()
}

func (cstr *Caster) distributeMessage(m string) {
	cstr.RLock()
	for c, _ := range cstr.Connections {
		c.outbound <- m
	}
	cstr.RUnlock()
}
