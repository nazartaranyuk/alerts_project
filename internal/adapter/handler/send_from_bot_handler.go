package handler

import (
	"nazartaraniuk/alertsProject/internal/domain"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func SendFromBotHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "application/json")
		var message domain.TelegramBotMessage

		if err := c.Bind(&message); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "invalid JSON",
			})
		}

		response := map[string]any{
			"status": "ok",
		}

		logrus.Printf("Message from bot: %v", message)
		return c.JSON(http.StatusOK, response)
	}
}
