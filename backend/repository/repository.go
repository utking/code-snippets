package repository

import (
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"

	"xorm.io/xorm"
)

const (
	DefaultDBFilePath = "./snippets.db"
)

var (
	orm *xorm.Engine
	err error
)

func InitDB(dbFilePath string) {
	if dbFilePath == "" {
		orm, err = xorm.NewEngine("sqlite3", DefaultDBFilePath)
	} else {
		orm, err = xorm.NewEngine("sqlite3", dbFilePath)
	}
}

type CustomContext struct {
	echo.Context
}

func (c *CustomContext) DB() (*xorm.Engine, error) {
	return orm, err
}
