package handler

import (
	"encoding/json"
	"nazartaraniuk/alertsProject/internal/domain"
	"nazartaraniuk/alertsProject/internal/usecase"
	"net/http"
	"time"
)

// GetAlarms godoc
// @Summary Get all alarms
// @Description Returns all alarms info
// @Tags alarms
// @Produce json
// @Success 200 {array} domain.RegionAlarmInfo
// @Failure 500 {object} domain.Error
// @Router /alerts [get]
func GetAlarms(serv usecase.GetAlarmInfoService) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		response, err := serv.GetCurrentAlerts()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(
				domain.NewError("Server error, api might have limit, try again after 1 minute", time.Now()),
			)
			return
		}

		_ = json.NewEncoder(w).Encode(response)
	}
}
