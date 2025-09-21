package handler

import (
	"encoding/json"
	"nazartaraniuk/alertsProject/internal/domain"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// Health godoc
// @Summary check server health
// @Description checks server health and returns result
// @Tags health
// Success
func Health() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(c.Response().Writer).Encode(domain.Health{
			Status: "ok",
			Time:   time.Now().Format(time.RFC3339),
		})
		if err != nil {
			logrus.Printf("Cannot encode: %v", err)
		}
		return err
	}
}
