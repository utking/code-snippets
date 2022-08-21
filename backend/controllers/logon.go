package controllers

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/utking/code-snippets/middleware"
	"github.com/utking/code-snippets/repository"
	"github.com/utking/code-snippets/types"

	"github.com/labstack/echo/v4"
	"xorm.io/xorm"
)

func Login(c echo.Context) error {
	var err error

	if IsLoggedIn(c) {
		return c.Redirect(http.StatusSeeOther, "/")
	}

	if c.Request().Method == "POST" {
		params, errParams := c.FormParams()

		if errParams == nil {
			_username := strings.Trim(params.Get("username"), " ")
			_password := strings.Trim(params.Get("password"), " ")

			if _, err = ValidateUser(_username, _password, c); err == nil {
				return c.Redirect(http.StatusSeeOther, "/")
			}

			err = fmt.Errorf("wring credentials")
		}
	}

	return c.Render(http.StatusOK, "login.html", map[string]interface{}{
		"Error": err,
	})
}

func Register(c echo.Context) error {
	var err error

	if IsLoggedIn(c) {
		return c.Redirect(http.StatusSeeOther, "/")
	}

	if c.Request().Method == "POST" {
		params, errParams := c.FormParams()

		if errParams == nil {
			_username := strings.Trim(params.Get("username"), " ")
			_password := strings.Trim(params.Get("password"), " ")
			_confirmation := strings.Trim(params.Get("confirmation"), " ")

			if _password == "" || _password != _confirmation {
				err = fmt.Errorf("password and confirmation do not match")
			} else {
				if err = createUser(_username, _password, c); err == nil {
					return c.Redirect(http.StatusSeeOther, "/")
				}
			}
		}
	}

	return c.Render(http.StatusOK, "register.html", map[string]interface{}{
		"Error": err,
	})
}

func createUser(_username, _password string, c echo.Context) error {
	var (
		err   error
		db    *xorm.Engine
		count int64
	)

	reUsername := regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_-]{2,31}$")
	rePassword := regexp.MustCompile(`^\w{8,16}$`)

	if !reUsername.MatchString(_username) {
		return fmt.Errorf("username must be 3 to 32 characters and start with a letter")
	}

	if !rePassword.MatchString(_password) {
		return fmt.Errorf("password must be 8 to 16 non-space characters")
	}

	existingUser := new(types.User)
	user := new(types.User)

	existingUser.Username = _username

	user.Active = true
	user.Username = _username
	user.IsAdmin = false
	user.Hash, err = user.HashPassword(_password)

	if err != nil {
		return err
	}

	if cc, ok := c.(*repository.CustomContext); ok {
		if db, err = cc.DB(); err == nil {
			if count, _ = db.Where("username=?", existingUser.Username).Count(existingUser); count > 0 {
				return cc.JSON(http.StatusConflict, map[string]interface{}{
					"Error":  "Try another username",
					"Status": http.StatusConflict,
				})
			}

			if _, err = db.InsertOne(*user); err == nil {
				_, err = ValidateUser(_username, _password, c)
			}
		}
	}

	return err
}

func ValidateUser(_username, _password string, c echo.Context) (bool, error) {
	var (
		err error
		db  *xorm.Engine
	)

	if IsLoggedIn(c) {
		return true, nil
	}

	user := new(types.User)
	user.Username = _username
	err = fmt.Errorf("wrong credentials")

	if cc, ok := c.(*repository.CustomContext); ok {
		if db, err = cc.DB(); err == nil {
			if _, err = db.Get(user); err == nil {
				if user.CheckPasswordHash(_password, user.Hash) {
					if err = SaveUserToSession(c, _username, user.ID); err == nil {
						return true, nil
					}
				}
			}
		}
	}

	return false, err
}

func Logout(c echo.Context) error {
	return middleware.LogoutHandler(c, "/")
}
