package controllers

import (
	. "code-snippets/repository"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"xorm.io/xorm"
)

type Note struct {
	Title   string `xorm:"varchar(32) NOT NULL"`
	Content string `xorm:"TEXT NOT NULL"`
	TagID   uint16 `xorm:"tag_id NOT NULL"`
	ID      uint16 `xorm:"id pk autoincr"`
	Indent  uint8  `xorm:"NOT NULL"`
}

type Notes []Note

func GetNotes(c echo.Context) error {
	notes := make(Notes, 0)

	cc, ok := c.(*CustomContext)

	if ok {
		if db, err := cc.DB(); err == nil {
			_ = db.CreateTables(Note{})

			if err = db.Find(&notes); err != nil {
				fmt.Println(err)
			}
		}
	}

	return cc.JSON(http.StatusOK, notes)
}

func GetNote(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, strconv.IntSize)

	if err == nil {
		var (
			note Note
			has  bool
		)

		cc, ok := c.(*CustomContext)

		if ok {
			if db, err := cc.DB(); err == nil {
				_ = db.CreateTables(Note{})

				if has, err = db.ID(id).Get(&note); err == nil && has {
					return c.JSON(http.StatusOK, note)
				}
			}
		}
	}

	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"Error":  "Snippet not found",
		"Status": http.StatusNotFound,
	})
}

func GetCategoryNotes(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, strconv.IntSize)
	notes := make(Notes, 0)

	if err == nil {
		cc, ok := c.(*CustomContext)

		if ok {
			if db, err := cc.DB(); err == nil {
				_ = db.CreateTables(Note{})

				if err = db.Find(&notes, &Note{TagID: uint16(id)}); err != nil {
					fmt.Println(err)
				}
			}
		}
	}

	return c.JSON(http.StatusOK, notes)
}

func DeleteNote(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, strconv.IntSize)

	if err == nil {
		cc, ok := c.(*CustomContext)

		if ok {
			if db, err := cc.DB(); err == nil {
				_ = db.CreateTables(Note{})

				if _, err = db.ID(id).Delete(&Note{}); err == nil {
					return c.JSON(http.StatusOK, map[string]interface{}{
						"Status": http.StatusOK,
					})
				}

				fmt.Println(err)
			}
		}
	}

	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"Error":  "Snippet not found",
		"Status": http.StatusNotFound,
	})
}

func PostNote(c echo.Context) error {
	var (
		count, noteID int64
		result        sql.Result
	)

	note := new(Note)
	cc, ok := c.(*CustomContext)

	if err := cc.Bind(note); err != nil {
		return cc.JSON(http.StatusConflict, map[string]interface{}{
			"Error":  "Wrong parameters",
			"Status": http.StatusBadRequest,
		})
	}

	if note.TagID <= 0 {
		return cc.JSON(http.StatusConflict, map[string]interface{}{
			"Error":  "Wrong or missing tag",
			"Status": http.StatusBadRequest,
		})
	}

	if strings.Trim(note.Title, " ") == "" || strings.Trim(note.Content, " ") == "" {
		return cc.JSON(http.StatusConflict, map[string]interface{}{
			"Error":  "Empty title or content is not allowed",
			"Status": http.StatusBadRequest,
		})
	}

	if ok {
		if db, err := cc.DB(); err == nil {
			_ = db.CreateTables(Note{})

			if count, _ = db.Where("tag_id=? AND title=?", note.TagID, note.Title).Count(&Note{}); count > 0 {
				return cc.JSON(http.StatusConflict, map[string]interface{}{
					"Error":  "The snippet with this alias exists for the tag",
					"Status": http.StatusConflict,
				})
			}

			if _, err = db.InsertOne(*note); err != nil {
				return cc.JSON(http.StatusConflict, map[string]interface{}{
					"Error":  err.Error(),
					"Status": http.StatusNotFound,
				})
			}

			if result, err = db.Exec("SELECT last_insert_rowid()"); err != nil {
				return cc.JSON(http.StatusConflict, map[string]interface{}{
					"Error":  err.Error(),
					"Status": http.StatusNotFound,
				})
			}

			if noteID, err = result.LastInsertId(); err != nil {
				return cc.JSON(http.StatusConflict, map[string]interface{}{
					"Error":  err.Error(),
					"Status": http.StatusNotFound,
				})
			}

			var newNote Note

			if _, err = db.ID(noteID).Get(&newNote); err == nil {
				return c.JSON(http.StatusOK, newNote)
			}
		}
	}

	return cc.JSON(http.StatusConflict, map[string]interface{}{
		"Error":  "Bad context",
		"Status": http.StatusNotFound,
	})
}

func PutNote(c echo.Context) error {
	var (
		existingNote Note
		db           *xorm.Engine
	)

	if id, err := strconv.ParseInt(c.Param("id"), 10, strconv.IntSize); err == nil {
		newNote := new(Note)
		cc, ok := c.(*CustomContext)

		if ok {
			if err := cc.Bind(newNote); err != nil {
				return cc.JSON(http.StatusConflict, map[string]interface{}{
					"Error":  "Wrong parameters",
					"Status": http.StatusBadRequest,
				})
			}

			newNote.ID = uint16(id)

			if newNote.TagID <= 0 {
				return cc.JSON(http.StatusConflict, map[string]interface{}{
					"Error":  "Wrong tag",
					"Status": http.StatusBadRequest,
				})
			}

			if strings.Trim(newNote.Title, " ") == "" || strings.Trim(newNote.Content, " ") == "" {
				return cc.JSON(http.StatusConflict, map[string]interface{}{
					"Error":  "Empty title or content is not allowed",
					"Status": http.StatusBadRequest,
				})
			}

			if db, err = cc.DB(); err == nil {
				_ = db.CreateTables(Note{})

				if _, err = db.ID(id).Get(&existingNote); err != nil {
					return cc.JSON(http.StatusConflict, map[string]interface{}{
						"Error":  err.Error(),
						"Status": http.StatusNotFound,
					})
				}

				_, err = db.ID(id).Update(*newNote)

				if err != nil {
					return cc.JSON(http.StatusConflict, map[string]interface{}{
						"Error":  err.Error(),
						"Status": http.StatusConflict,
					})
				}

				var newNote Note

				if _, err = db.ID(id).Get(&newNote); err == nil {
					return c.JSON(http.StatusOK, newNote)
				}

			}

			return cc.JSON(http.StatusConflict, map[string]interface{}{
				"Error":  err.Error(),
				"Status": http.StatusNotFound,
			})
		}
	}

	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"Error":  "Bad context",
		"Status": http.StatusNotFound,
	})
}
