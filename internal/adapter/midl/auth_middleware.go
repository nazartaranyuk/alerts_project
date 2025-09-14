package midl

import (
	"crypto/subtle"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func AddTestAuthMiddleWare(server *echo.Echo, testUsername string, testPassword string) {
	server.Use(middleware.BasicAuth(func(username, password string, _ echo.Context) (bool, error) {
		if subtle.ConstantTimeCompare([]byte(username), []byte(testUsername)) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte(testPassword)) == 1 {
			return true, nil
		}
		return false, nil
	}))
}
