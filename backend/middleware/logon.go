package middleware

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

var Store *sessions.CookieStore

func init() {
	Store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_TOKEN")))
}

func LogoutHandler(c echo.Context, redirectURL string) error {
	session, _ := Store.Get(c.Request(), "session.id")
	session.Values["authenticated"] = false
	session.Values["username"] = ""
	session.Values["id"] = nil
	_ = session.Save(c.Request(), c.Response().Writer)

	return c.Redirect(http.StatusSeeOther, redirectURL)
}
