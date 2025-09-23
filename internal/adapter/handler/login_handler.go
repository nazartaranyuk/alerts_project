package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"nazartaraniuk/alertsProject/internal/config"
	"nazartaraniuk/alertsProject/internal/domain"
	"nazartaraniuk/alertsProject/internal/usecase"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// LoginHandler godoc
// @Summary      Login
// @Description  Authorize user and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        loginReq  body      domain.LoginReq  true  "Login credentials"
// @Success      200       {object}  domain.TokenResp
// @Failure      400       {string}  string "Bad Request"
// @Failure      401       {string}  string "Unauthorized"
// @Failure      500       {string}  string "Internal Server Error"
// @Router       /login [post]
func LoginHandler(cfg config.Config, service usecase.UserService) echo.HandlerFunc {
	return func(c echo.Context) error {
		var loginReq domain.LoginReq

		if err := json.NewDecoder(c.Request().Body).Decode(&loginReq); err != nil {
			return c.NoContent(http.StatusBadRequest)
		}
		u, err := service.LoginUser(loginReq)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return c.String(http.StatusUnauthorized, "unauthorized")
			}
			return c.NoContent(http.StatusInternalServerError)
		}
		if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(loginReq.Password)); err != nil {
			return c.String(http.StatusUnauthorized, "invalid credentials")
		}

		ttl := 15 * time.Minute
		now := time.Now()
		claims := jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(u.ID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signed, err := token.SignedString([]byte(cfg.Server.JWTSecret))
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, domain.TokenResp{
			AccessToken: signed,
			ExpiresAt:   claims.ExpiresAt.Unix(),
		})
	}
}
