package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"code-snippets/middleware"

	"github.com/labstack/echo/v4"
)

var (
	username, password string
)

func InitUser(_username, _password string) error {
	if strings.Trim(_username, " ") == "" || strings.Trim(_password, " ") == "" {
		return fmt.Errorf("username and password cannot be empty")
	}

	username = _username
	password = _password

	return nil
}

func Login(c echo.Context) error {
	var err error

	if c.Request().Method == "POST" {
		params, errParams := c.FormParams()

		if errParams != nil {
			log.Fatal(errParams)
		}

		username := params.Get("username")
		password := params.Get("password")

		if _, err = ValidateUser(username, password, c); err == nil {
			return c.Redirect(http.StatusSeeOther, "/")
		}

		err = fmt.Errorf("%s", err)
	}

	return c.Render(http.StatusOK, "login.html", map[string]interface{}{
		"Error": err,
	})
}

func ValidateUser(_username, _password string, c echo.Context) (bool, error) {
	session, _ := middleware.Store.Get(c.Request(), "session.id")

	if (session.Values["authenticated"] != nil) && session.Values["authenticated"] != false {
		return true, nil
	}

	if _username == username && _password == _password {
		session.Values["authenticated"] = true
		session.Values["username"] = username
		session.Save(c.Request(), c.Response().Writer)

		return true, nil
	}

	return false, fmt.Errorf("wrong credentials")
}

func Logout(c echo.Context) error {
	return middleware.LogoutHandler(c, "/")
}
