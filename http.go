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
			ws:       ws,
			outbound: make(chan string, 256),
		}

		cstr.Create <- c
		defer func() { cstr.Destroy <- c }()

		go c.write()

		c.read()
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("index").Parse(indexTemplate)
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, nil)
}

func InitHttp(listen *string, c *Caster) {
	http.HandleFunc("/", homeHandler)
	http.Handle("/ws", websocket.Handler(createWebsocketHandler(c)))

	log.Printf("listening to http on %s", *listen)
	go http.ListenAndServe(*listen, nil)
}
