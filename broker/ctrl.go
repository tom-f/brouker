package broker

import (
	"net/http"
)

// Control blah
type Ctrl struct {
	connections map[*Connection]bool
	broadcast   chan []byte
	register    chan *Connection
	unregister  chan *Connection
}

func NewCtrl() *Ctrl {
	return &Ctrl{
		broadcast:   make(chan []byte),
		register:    make(chan *Connection),
		unregister:  make(chan *Connection),
		connections: make(map[*Connection]bool),
	}
}

func (h *Ctrl) Run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
			}
		case m := <-h.broadcast:
			for c := range h.connections {
				select {
				case c.send <- m:
				default:
					delete(h.connections, c)
					close(c.send)
				}
			}
		}
	}
}

type CtrlWriter struct {
	C *Ctrl
}

func (wsW CtrlWriter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := make([]byte, r.ContentLength)
	_, err := r.Body.Read(p)
	if err != nil {

	}
	wsW.C.broadcast <- p
}
