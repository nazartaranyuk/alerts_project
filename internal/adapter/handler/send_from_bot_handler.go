package handler

import (
	"nazartaraniuk/alertsProject/internal/domain"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// SendFromBotHandler godoc
// @Summary      Receive message from bot (Experimental)
// @Description  Accepts a message from the Telegram bot and logs it. This endpoint is experimental and may change.
// @Tags         bot
// @Accept       json
// @Produce      json
// @Param        message  body      domain.TelegramBotMessage  true  "Telegram Bot Message"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]string
// @Router       /bot/send [post]
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
