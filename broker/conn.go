package broker

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// Connection blah blah
type Connection struct {
	ws   *websocket.Conn
	send chan []byte
	h    *Ctrl
}

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

func (c *Connection) reader() {
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		c.h.broadcast <- message
	}
	c.ws.Close()
}

func (c *Connection) writer() {
	for message := range c.send {
		err := c.ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
	c.ws.Close()
}

// ConnHandler handles all connections to the supplied control
type ConnHandler struct {
	C *Ctrl
}

func (wsh ConnHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	c := &Connection{send: make(chan []byte, 256), ws: ws, h: wsh.C}
	c.h.register <- c
	defer func() { c.h.unregister <- c }()
	go c.writer()
	c.reader()
}
