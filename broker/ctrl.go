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

func (ctrl *Ctrl) Run() {
	for {
		select {
		case conn := <-ctrl.register:
			ctrl.connections[conn] = true
		case conn := <-ctrl.unregister:
			if _, ok := ctrl.connections[conn]; ok {
				delete(ctrl.connections, conn)
				close(conn.send)
			}
		case m := <-ctrl.broadcast:
			for conn := range ctrl.connections {
				select {
				case conn.send <- m:
				default:
					delete(ctrl.connections, conn)
					close(conn.send)
				}
			}
		}
	}
}

type CtrlWriter struct {
	C *Ctrl
}

func (ctrlW CtrlWriter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := make([]byte, r.ContentLength)
	_, err := r.Body.Read(p)
	if err != nil {

	}
	ctrlW.C.broadcast <- p
}
