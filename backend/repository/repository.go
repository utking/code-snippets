package repository

import (
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

var (
	orm *xorm.Engine
	err error
)

func init() {
	orm, err = xorm.NewEngine("sqlite3", "./snippets.db")
}

type CustomContext struct {
	echo.Context
}

func (c *CustomContext) DB() (*xorm.Engine, error) {
	return orm, err
}
