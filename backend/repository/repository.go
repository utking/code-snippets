package repository

import (
	. "code-snippets/types"

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

func InitDB(dbFilePath string) (*xorm.Engine, error) {
	var _dbFilePath string

	if dbFilePath == "" {
		_dbFilePath = DefaultDBFilePath
	} else {
		_dbFilePath = dbFilePath
	}

	if orm, err = xorm.NewEngine("sqlite3", _dbFilePath); err != nil {
		return nil, err
	}

	return orm, nil
}

func InitTables() error {
	if tagError := orm.CreateTables(NoteTag{}); err != nil {
		return tagError
	}

	if noteError := orm.CreateTables(Note{}); err != nil {
		return noteError
	}

	return nil
}

type CustomContext struct {
	echo.Context
}

func (c *CustomContext) DB() (*xorm.Engine, error) {
	return orm, err
}
