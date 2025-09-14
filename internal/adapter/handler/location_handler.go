package handler

import (
	"log"
	"nazartaraniuk/alertsProject/internal/adapter/ws"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type message struct {
	Type     string  `json:"type"`
	UserID   string  `json:"user_id"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Accuracy float64 `json:"accuracy"`
	TS       string  `json:"ts"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin:     func(_ *http.Request) bool { return true },
}

func Handler(h *ws.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "user_id required", http.StatusBadRequest)
			return
		}
		wb, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		c := ws.NewConn(wb)
		h.Join(userID, c)
		go c.WriteLoop()
		defer func() {
			h.Leave(userID, c)
			c.Close()
		}()
		for {
			var m message
			if err := c.ReadJSON(&m); err != nil {
				return
			}
			if m.TS == "" {
				m.TS = time.Now().UTC().Format(time.RFC3339)
			}
			if m.Type == "loc" && m.UserID == userID {
				log.Printf("[WS] user %s -> lat=%.6f lon=%.6f accuracy=%.1f ts=%s",
					m.UserID, m.Lat, m.Lon, m.Accuracy, m.TS)
				h.Broadcast(userID, ws.Payload{
					UserID:   m.UserID,
					Lat:      m.Lat,
					Lon:      m.Lon,
					Accuracy: m.Accuracy,
					TS:       m.TS,
				})
			}
		}
	}
}
