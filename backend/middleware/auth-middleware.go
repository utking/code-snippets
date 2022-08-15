package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	// AuthConfig defines the config for Auth middleware.
	AuthConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper

		// Validator is a function to validate Auth credentials.
		// Required.
		Validator AuthValidator

		// Redirect URL for authentication
		RedirectURL string
	}

	// AuthValidator defines a function to validate Auth credentials.
	AuthValidator func(string, string, echo.Context) (bool, error)
)

var (
	// DefaultAuthConfig is the default Auth middleware config.
	DefaultAuthConfig = AuthConfig{
		Skipper:     middleware.DefaultSkipper,
		RedirectURL: "/login",
	}
)

// Auth returns an Auth middleware.
//
// For valid credentials it calls the next handler.
// For missing or invalid credentials, it sends "401 - Unauthorized" response.
func Auth(fn AuthValidator) echo.MiddlewareFunc {
	c := DefaultAuthConfig
	c.Validator = fn

	return AuthWithConfig(c)
}

// AuthWithConfig returns an Auth middleware with config.
// See `Auth()`.
func AuthWithConfig(config AuthConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Validator == nil {
		panic("echo: auth middleware requires a validator function")
	}

	if config.Skipper == nil {
		config.Skipper = DefaultAuthConfig.Skipper
	}

	if config.RedirectURL == "" {
		config.RedirectURL = DefaultAuthConfig.RedirectURL
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			session, _ := Store.Get(c.Request(), "session.id")
			if (session.Values["authenticated"] != nil) && session.Values["authenticated"] != false {
				return next(c)
			} else {
				return c.Redirect(http.StatusSeeOther, config.RedirectURL)
			}

		}
	}
}
