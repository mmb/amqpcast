package amqpcast

import (
	"html/template"
	"log"
	"net/http"

	"code.google.com/p/go.net/websocket"
)

func createWebsocketHandler(cstr *Caster) func(ws *websocket.Conn) {
	return func(ws *websocket.Conn) {
		c := &Connection{
			Ws:       ws,
			Outbound: make(chan string, 256),
		}

		cstr.Create <- c
		defer func() { cstr.Destroy <- c }()

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

func InitHttp(c *Caster) {
	http.HandleFunc("/", homeHandler)
	http.Handle("/ws", websocket.Handler(createWebsocketHandler(c)))

	log.Println("listening to http on :12345")
	go http.ListenAndServe(":12345", nil)
}
