package amqpcast

type Caster struct {
	Connections map[*Connection]bool
	Create      chan *Connection
	Destroy     chan *Connection
	Outbound    chan string
}
