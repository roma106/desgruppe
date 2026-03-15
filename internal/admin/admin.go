package admin

import (
	"desgruppe/internal/config"
	"desgruppe/internal/logger"
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type Admin struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func LoginAdmin(c echo.Context, cfg *config.Config) error {
	admin := Admin{}
	login := cfg.AdminLogin
	password := cfg.AdminPassword
	// login := "daniil"
	// password := "031128"

	if err := json.NewDecoder(c.Request().Body).Decode(&admin); err != nil {
		logger.Error("Error parsing request body to log admin:", err)
		return c.String(http.StatusBadRequest, "Invalid request payload")
	}

	logger.Info("Admin tried to login: ", admin)

	if admin.Login == login && admin.Password == password {
		authcookie := new(http.Cookie)
		authcookie.Name = "admin-logged"
		authcookie.Value = "true"
		authcookie.Path = "/"
		authcookie.Expires = time.Now().Add(time.Hour * 2)
		// authcookie.Secure = true
		// authcookie.HttpOnly = true
		// authcookie.SameSite = http.SameSiteStrictMode

		c.SetCookie(authcookie)

		logger.Info("Welcome, ", admin.Login)

		return c.String(http.StatusOK, "Login successful")
	}

	return c.String(http.StatusForbidden, "Forbidden")
}

func AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authcookie, err := c.Cookie("admin-logged")

		if err != nil || authcookie.Value != "true" {
			return c.Redirect(http.StatusSeeOther, "/admin-login")
		}

		return next(c)
	}
}
