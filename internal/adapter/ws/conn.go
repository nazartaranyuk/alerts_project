package ws

import (
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

type Conn struct {
	ws  *websocket.Conn
	out chan []byte
}

func NewConn(ws *websocket.Conn) *Conn {
	return &Conn{ws: ws, out: make(chan []byte, 64)}
}

func (c *Conn) ReadJSON(v any) error {
	c.ws.SetReadLimit(1 << 16)
	err := c.ws.SetReadDeadline(time.Now().Add(60 * time.Second))
	if err != nil {
		logrus.Printf("Cannot set connection deadline: %v", err)
	}
	c.ws.SetPongHandler(func(string) error {
		err := c.ws.SetReadDeadline(time.Now().Add(60 * time.Second))
		if err != nil {
			return err
		}
		return nil
	})
	return c.ws.ReadJSON(v)
}

func (c *Conn) WriteLoop() {
	t := time.NewTicker(30 * time.Second)
	defer t.Stop()
	for {
		select {
		case b, ok := <-c.out:
			err := c.ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err != nil {
				return
			}
			if !ok {
				_ = c.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.ws.WriteMessage(websocket.TextMessage, b); err != nil {
				return
			}
		case <-t.C:
			err := c.ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err != nil {
				return
			}
			if err := c.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Conn) Send(p Payload) {
	b, _ := json.Marshal(p)
	select {
	case c.out <- b:
	default:
	}
}

func (c *Conn) Close() {
	close(c.out)
	_ = c.ws.Close()
}
