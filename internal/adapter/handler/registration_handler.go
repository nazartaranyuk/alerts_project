package handler

import (
	"encoding/json"
	"nazartaraniuk/alertsProject/internal/domain"
	"nazartaraniuk/alertsProject/internal/repository"
	"nazartaraniuk/alertsProject/internal/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

// RegistrationHandler godoc
// @Summary      Register a new user
// @Description  Registers a user with email, username, and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        registerReq  body      domain.RegisterReq  true  "Registration request"
// @Success      201          {object}  map[string]interface{}
// @Failure      400          {string}  string  "Bad Request"
// @Failure      409          {string}  string  "Email already exists"
// @Failure      500          {string}  string  "Internal Server Error"
// @Router       /register [post]
func RegistrationHandler(userService *usecase.UserService) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req domain.RegisterReq
		if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
			return c.NoContent(http.StatusBadRequest)
		}
		if req.Email == "" || req.Password == "" {
			return c.NoContent(http.StatusBadRequest)
		}
		id, err := userService.RegisterUser(req)
		if err != nil {
			switch err {
			case repository.ErrEmailAlreadyExists:
				return c.String(http.StatusConflict, "email already exists")
			case repository.ErrCannotCreateUser:
				return c.NoContent(http.StatusInternalServerError)
			}
		}

		resp := domain.UserCreatedResponse{
			ID:       id,
			Email:    req.Email,
			Username: req.Username,
		}
		return c.JSON(http.StatusCreated, resp)
	}

}
