package handler

import (
	"nazartaraniuk/alertsProject/internal/adapter/ws"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

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

func Handler(h *ws.Hub) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Request().URL.Query().Get("user_id")
		if userID == "" {
			http.Error(c.Response(), "user_id required", http.StatusBadRequest)
			return nil
		}
		wb, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			logrus.Println(err)
			return nil
		}
		conn := ws.NewConn(wb)
		h.Join(userID, conn)
		go conn.WriteLoop()
		defer func() {
			h.Leave(userID, conn)
			conn.Close()
		}()
		for {
			var m message
			if err := conn.ReadJSON(&m); err != nil {
				return nil
			}
			if m.TS == "" {
				m.TS = time.Now().UTC().Format(time.RFC3339)
			}
			if m.Type == "loc" && m.UserID == userID {
				logrus.Printf("[WS] user %s -> lat=%.6f lon=%.6f accuracy=%.1f ts=%s",
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
