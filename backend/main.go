package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/utking/code-snippets/config"
	"github.com/utking/code-snippets/controllers"
	auth "github.com/utking/code-snippets/middleware"
	"github.com/utking/code-snippets/repository"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	httpPort = 8080
)

func main() {
	var (
		NoRegistration bool
		port           int
		dbFilePath     string
		err            error
	)

	flag.IntVar(&port, "port", httpPort, "HTTP server port")
	flag.StringVar(&dbFilePath, "db", repository.DefaultDBFilePath, "SQLite3 DB file path")
	flag.BoolVar(&NoRegistration, "no-reg", false, "Disable self-registration")
	flag.Parse()

	if _, err = repository.InitDB(dbFilePath); err != nil {
		log.Fatal(err)
	}
	if err = repository.InitTables(); err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	if err = config.InitTemplates(e); err != nil {
		log.Fatal(err)
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{
			echo.HeaderAccessControlAllowOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAccessControlAllowCredentials,
		},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) (bool, error) {
			return true, nil
		},
	}))

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &repository.CustomContext{c}
			return next(cc)
		}
	})

	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "cookie:_csrf",
	}))

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(auth.AuthWithConfig(auth.AuthConfig{
		Skipper: func(c echo.Context) bool {
			// For shared snippets
			if strings.HasPrefix(c.Request().URL.Path, "/fe/assets/") ||
				strings.HasPrefix(c.Request().URL.Path, "/note/s/") ||
				strings.HasPrefix(c.Request().URL.Path, "/note/fe/") {
				return true
			}

			// For login and register pages
			for _, url := range []string{
				"/login", "/register",
			} {
				if c.Request().URL.Path == url {
					return true
				}
			}

			return false
		},
		Validator: controllers.ValidateUser,
	}))
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{Root: "./views"}))

	e.GET("/", controllers.Loader)
	e.Match([]string{"GET", "POST"}, "/login", controllers.Login)

	if !NoRegistration {
		e.Match([]string{"GET", "POST"}, "/register", controllers.Register)
	} else {
		e.GET("/register", controllers.Login)
	}

	e.POST("/logout", controllers.Logout)

	tagGroup := e.Group("/tag")
	tagGroup.GET("", controllers.GetTags)
	tagGroup.GET("/", controllers.GetTags)
	tagGroup.PUT("/:tag", controllers.PutTag)

	noteGroup := e.Group("/note")
	// REST-like
	noteGroup.GET("/:id", controllers.GetNote)
	noteGroup.POST("/", controllers.PostNote)
	noteGroup.DELETE("/:id", controllers.DeleteNote)
	noteGroup.PUT("/:id", controllers.PutNote)
	// Custom
	noteGroup.GET("/tag/:tag", controllers.GetTagNotes)
	noteGroup.GET("/s/:hash", controllers.GetSharedNote)
	noteGroup.GET("/fe/:hash", controllers.Loader)
	noteGroup.POST("/share", controllers.ShareNote)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}
