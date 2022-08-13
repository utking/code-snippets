package controllers

import (
	. "code-snippets/repository"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

type NoteTag struct {
	Alias string `xorm:"varchar(32) not null"`
	ID    uint16 `xorm:"id pk autoincr"`
}

func GetTags(c echo.Context) error {
	var tags []NoteTag

	cc, ok := c.(*CustomContext)

	if ok {
		if db, err := cc.DB(); err == nil {
			_ = db.CreateTables(NoteTag{})
			_ = db.SQL("SELECT * from note_tag").Find(&tags)
		}
	}

	return cc.JSON(http.StatusOK, tags)
}

func PostTag(c echo.Context) error {
	var (
		count  int64
		result int64
	)

	tag := new(NoteTag)
	cc, ok := c.(*CustomContext)

	if err := cc.Bind(tag); err != nil {
		return cc.JSON(http.StatusConflict, map[string]interface{}{
			"Error":  "Wrong parameters",
			"Status": http.StatusBadRequest,
		})
	}

	if strings.Trim(tag.Alias, " ") == "" {
		return cc.JSON(http.StatusConflict, map[string]interface{}{
			"Error":  "Empty alias is not allowed",
			"Status": http.StatusBadRequest,
		})
	}

	if ok {
		if db, err := cc.DB(); err == nil {
			_ = db.CreateTables(NoteTag{})

			count, _ = db.Where("alias=?", tag.Alias).Count(&NoteTag{})
			if count > 0 {
				return cc.JSON(http.StatusConflict, map[string]interface{}{
					"Error":  "Tag exists",
					"Status": http.StatusConflict,
				})
			}

			result, err = db.InsertOne(*tag)

			if err != nil {
				return cc.JSON(http.StatusConflict, map[string]interface{}{
					"Error":  err.Error(),
					"Status": http.StatusNotFound,
				})
			}

			fmt.Println(result)
		}
	}

	return cc.JSON(http.StatusOK, map[string]interface{}{
		"Status": http.StatusOK,
	})
}

func GetTag(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, strconv.IntSize)

	if err == nil {
		var tags []NoteTag

		cc, ok := c.(*CustomContext)

		if ok {
			if db, err := cc.DB(); err == nil {
				_ = db.CreateTables(NoteTag{})
				err = db.SQL("SELECT * FROM note_tag WHERE id=?", id).Find(&tags)
				fmt.Println(err)

				if err == nil && len(tags) > 0 {
					return c.JSON(http.StatusOK, tags[len(tags)-1])
				}
			}
		}
	}

	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"Error":  "Tag not found",
		"Status": http.StatusNotFound,
	})
}

func DeleteTag(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, strconv.IntSize)

	if err == nil {
		cc, ok := c.(*CustomContext)

		if ok {
			if db, err := cc.DB(); err == nil {
				_ = db.CreateTables(NoteTag{})
				_, err = db.ID(id).Delete(&NoteTag{})
				fmt.Println(err)

				if err == nil {
					return c.JSON(http.StatusOK, map[string]interface{}{
						"Status": http.StatusOK,
					})
				}
			}
		}
	}

	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"Error":  "Tag not found",
		"Status": http.StatusNotFound,
	})
}
