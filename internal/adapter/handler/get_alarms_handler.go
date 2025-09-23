package handler

import (
	"encoding/json"
	"nazartaraniuk/alertsProject/internal/domain"
	"nazartaraniuk/alertsProject/internal/usecase"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	ServerErrorMessage = "Server error, api might have limit, try again after 1 minute"
)

// GetAlarms godoc
// @Summary Get all alarms
// @Description Returns all alarms info
// @Tags alarms
// @Produce json
// @Success 200 {array} domain.RegionAlarmInfo
// @Failure 500 {object} domain.Error
// @Router /alerts [get]
func GetAlarms(serv *usecase.GetAlarmInfoService) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Request().Header.Set("Content-Type", "application/json; charset=utf-8")

		response, err := serv.GetCurrentAlerts()
		if err != nil {
			c.Response().WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(c.Response().Writer).Encode(
				domain.NewError(ServerErrorMessage, time.Now()),
			)
		}

		_ = json.NewEncoder(c.Response().Writer).Encode(response)
		return err
	}
}
