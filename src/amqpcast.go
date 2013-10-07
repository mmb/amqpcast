package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"

	"code.google.com/p/go.net/websocket"
	"github.com/streadway/amqp"
)

type connection struct {
	ws       *websocket.Conn
	outbound chan string
}

type caster struct {
	connections map[*connection]bool
	create      chan *connection
	destroy     chan *connection
	outbound    chan string
}

func handleDeliveries(deliveries <-chan amqp.Delivery, cstr *caster) {
	for d := range deliveries {
		log.Printf("received AMQP message (%dB): %q", len(d.Body), d.Body)
		cstr.outbound <- string(d.Body[:])
	}
}

func initAmqp(amqpUrl *string, amqpExchange *string, amqpKey *string, cstr *caster) {
	amqpConn, err := amqp.Dial(*amqpUrl)
	if err != nil {
		log.Fatal(err)
	}

	amqpChannel, err := amqpConn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	queue, err := amqpChannel.QueueDeclare(
		"",    // name of the queue
		false, // durable
		true,  // autodelete
		true,  // exclusive
		false, // nowait
		nil,   // args
	)
	if err != nil {
		log.Fatal(err)
	}

	amqpChannel.QueueBind(
		queue.Name,
		*amqpKey,
		*amqpExchange,
		false, // nowait
		nil,   // args
	)

	deliveries, err := amqpChannel.Consume(
		queue.Name,
		"",    // consumer
		false, // autoack
		false, // exclusive
		false, // nolocal
		false, // nowait
		nil,   // args
	)

	go handleDeliveries(deliveries, cstr)
}

func (c *connection) write() {
	for message := range c.outbound {
		err := websocket.Message.Send(c.ws, message)
		if err != nil {
			break
		}
	}
}

func (c *connection) read() {
	for {
		var message string
		err := websocket.Message.Receive(c.ws, &message)
		if err != nil {
			break
		}
	}
}

func createWebsocketHandler(cstr *caster) func(ws *websocket.Conn) {
	return func(ws *websocket.Conn) {
		c := &connection{
			ws:       ws,
			outbound: make(chan string, 256),
		}

		cstr.create <- c
		defer func() { cstr.destroy <- c }()

		go c.write()

		c.read()
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, nil)
}

func initHttp(c *caster) {
	http.HandleFunc("/", homeHandler)
	http.Handle("/ws", websocket.Handler(createWebsocketHandler(c)))

	go http.ListenAndServe(":12345", nil)
}

func main() {
	var amqpUrl = flag.String("amqp-url", "amqp://localhost:5672/",
		"AMQP server URL")
	var amqpExchange = flag.String("amqp-exchange", "test", "AMQP exchange")
	var amqpKey = flag.String("amqp-key", "test", "AMQP routing key")

	flag.Parse()

	var cstr = caster{
		connections: make(map[*connection]bool),
		create:      make(chan *connection),
		destroy:     make(chan *connection),
		outbound:    make(chan string),
	}

	initHttp(&cstr)
	initAmqp(amqpUrl, amqpExchange, amqpKey, &cstr)

	for {
		select {
		case c := <-cstr.create:
			log.Printf("new client")
			cstr.connections[c] = true
		case c := <-cstr.destroy:
			log.Printf("client closed")
			delete(cstr.connections, c)
			c.ws.Close()
			close(c.outbound)
		case m := <-cstr.outbound:
			for c, _ := range cstr.connections {
				c.outbound <- m
			}
		}
	}
}
