package handler

import (
	"encoding/json"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func Health() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(c.Response().Writer).Encode(map[string]string{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
		if err != nil {
			logrus.Printf("Cannot encode: %v", err)
		}
		return err
	}
}
