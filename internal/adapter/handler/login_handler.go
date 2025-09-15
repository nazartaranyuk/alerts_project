package handler

import (
	"nazartaraniuk/alertsProject/internal/config"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// LoginHandler godoc
// @Summary Login
// @Description authorise user with JWT token
// @Tags alarms
// @Produce json
// @Success 200 {map} success
// @Failure 500 {object} 505
// @Router /alerts [get]
func LoginHandler(cfg config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		JWTSecret := []byte(cfg.Server.JWTSecret)

		type Credentials struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		var creds Credentials
		if err := c.Bind(&creds); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
		}

		if creds.Email != cfg.Server.AdminEmail || creds.Password != cfg.Server.AdminPassword {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		}

		claims := jwt.MapClaims{
			"sub":   "user-1",
			"email": creds.Email,
			"role":  "user",
			"exp":   time.Now().Add(24 * time.Hour).Unix(),
			"iat":   time.Now().Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		signed, err := token.SignedString(JWTSecret)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "cannot sign token"})
		}

		return c.JSON(http.StatusOK, map[string]any{
			"access_token": signed,
			"token_type":   "Bearer",
			"expires_in":   24 * 3600,
		})
	}
}
