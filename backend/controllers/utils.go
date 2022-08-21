package controllers

import (
	"fmt"

	"github.com/utking/code-snippets/middleware"

	"github.com/labstack/echo/v4"
)

func GetUserID(c echo.Context) (uint16, error) {
	session, err := middleware.Store.Get(c.Request(), "session.id")

	if err != nil {
		return 0, err
	}

	if (session.Values["authenticated"] == nil) || session.Values["authenticated"] == false {
		return 0, fmt.Errorf("the user is not authenticated")
	}

	if _userID, ok := session.Values["id"].(uint16); ok {
		return _userID, nil
	}

	return 0, fmt.Errorf("wrong user ID. Try to log out and log in again")
}

func IsLoggedIn(c echo.Context) bool {
	session, err := middleware.Store.Get(c.Request(), "session.id")

	if err != nil {
		return false
	}

	return (session.Values["authenticated"] != nil) && session.Values["authenticated"] != false
}

func SaveUserToSession(c echo.Context, _username string, _userID uint16) error {
	session, err := middleware.Store.Get(c.Request(), "session.id")

	if err != nil {
		return err
	}

	session.Values["authenticated"] = true
	session.Values["username"] = _username
	session.Values["id"] = _userID
	err = session.Save(c.Request(), c.Response().Writer)

	if err != nil {
		return err
	}

	return nil
}

func SessionFlash(c echo.Context, _data interface{}) error {
	session, err := middleware.Store.Get(c.Request(), "session.id")

	if err != nil {
		return err
	}

	session.Values["flash"] = _data
	err = session.Save(c.Request(), c.Response().Writer)

	if err != nil {
		return err
	}

	return nil
}
