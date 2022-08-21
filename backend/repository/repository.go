package repository

import (
	"github.com/utking/code-snippets/types"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"

	"xorm.io/xorm"
)

const (
	DefaultDBFilePath = "./snippets.db"
)

var (
	orm  *xorm.Engine
	_err error
)

func InitDB(dbFilePath string) (*xorm.Engine, error) {
	var (
		_dbFilePath string
	)

	if dbFilePath == "" {
		_dbFilePath = DefaultDBFilePath
	} else {
		_dbFilePath = dbFilePath
	}

	if orm, _err = xorm.NewEngine("sqlite3", _dbFilePath); _err != nil {
		return nil, _err
	}

	return orm, nil
}

func InitTables() error {
	if exists, _ := orm.IsTableExist(types.User{}); !exists {
		if err := orm.CreateTables(types.User{}); err != nil {
			return err
		}

		if err := orm.CreateIndexes(types.User{}); err != nil {
			return err
		}

		if err := orm.CreateUniques(types.User{}); err != nil {
			return err
		}
	}

	if exists, _ := orm.IsTableExist(types.Note{}); !exists {
		if err := orm.CreateTables(types.Note{}); err != nil {
			return err
		}

		if err := orm.CreateIndexes(types.Note{}); err != nil {
			return err
		}

		if err := orm.CreateUniques(types.Note{}); err != nil {
			return err
		}
	}

	if exists, _ := orm.IsTableExist(types.SharedNote{}); !exists {
		if err := orm.CreateTables(types.SharedNote{}); err != nil {
			return err
		}

		if err := orm.CreateIndexes(types.SharedNote{}); err != nil {
			return err
		}

		if err := orm.CreateUniques(types.SharedNote{}); err != nil {
			return err
		}
	}

	return nil
}

type CustomContext struct {
	echo.Context
}

func (c *CustomContext) DB() (*xorm.Engine, error) {
	return orm, _err
}
