package main

import (
	"code-snippets/config"
	"code-snippets/controllers"
	auth "code-snippets/middleware"
	"code-snippets/repository"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	httpPort = 8080
)

func main() {
	var (
		port       int
		dbFilePath string
		err        error
	)

	flag.IntVar(&port, "port", httpPort, "HTTP server port")
	flag.StringVar(&dbFilePath, "db", repository.DefaultDBFilePath, "SQLite3 DB file path")
	flag.Parse()

	repository.InitDB(dbFilePath)

	if err = controllers.InitUser(os.Getenv("CS_USER"), os.Getenv("CS_PASSWORD")); err != nil {
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
			return c.Request().URL.Path == "/login"
		},
		Validator: controllers.ValidateUser,
	}))
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{Root: "./views"}))

	e.GET("/", controllers.Loader)
	e.Match([]string{"GET", "POST"}, "/login", controllers.Login)
	e.POST("/logout", controllers.Logout)

	tagGroup := e.Group("/tag")
	tagGroup.GET("", controllers.GetTags)
	tagGroup.GET("/", controllers.GetTags)
	tagGroup.POST("/", controllers.PostTag)
	tagGroup.GET("/:id", controllers.GetTag)
	tagGroup.DELETE("/:id", controllers.DeleteTag)

	noteGroup := e.Group("/note")
	noteGroup.GET("", controllers.GetNotes)
	noteGroup.GET("/", controllers.GetNotes)
	noteGroup.GET("/:id", controllers.GetNote)
	noteGroup.POST("/", controllers.PostNote)
	noteGroup.DELETE("/:id", controllers.DeleteNote)
	noteGroup.PUT("/:id", controllers.PutNote)
	noteGroup.GET("/category/:id", controllers.GetCategoryNotes)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}
