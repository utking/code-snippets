package controllers

import (
	. "code-snippets/repository"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type Note struct {
	Title   string `xorm:"varchar(32) NOT NULL"`
	Syntax  string `xorm:"varchar(32) NOT NULL"`
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
				has, err = db.ID(id).Get(&note)

				if err == nil && has {
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
				err = db.Find(&notes, &Note{TagID: uint16(id)})
				fmt.Println(err)
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
				_, err = db.ID(id).Delete(&Note{})
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
		"Error":  "Snippet not found",
		"Status": http.StatusNotFound,
	})
}

func PostNote(c echo.Context) error {
	var (
		count int64
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
			"Error":  "Wrong tag",
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

			count, _ = db.Where("tag_id=? AND title=?", note.TagID, note.Title).Count(&Note{})
			if count > 0 {
				return cc.JSON(http.StatusConflict, map[string]interface{}{
					"Error":  "The snippet with this alias exists for the tag",
					"Status": http.StatusConflict,
				})
			}

			_, err = db.InsertOne(*note)

			if err != nil {
				return cc.JSON(http.StatusConflict, map[string]interface{}{
					"Error":  err.Error(),
					"Status": http.StatusNotFound,
				})
			}
		}
	}

	return cc.JSON(http.StatusOK, map[string]interface{}{
		"Status": http.StatusOK,
	})
}

func PutNote(c echo.Context) error {
	var (
		has          bool
		existingNote Note
	)

	id, err := strconv.ParseInt(c.Param("id"), 10, strconv.IntSize)

	if err == nil {
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

			if db, err := cc.DB(); err == nil {
				_ = db.CreateTables(Note{})

				if has, err = db.ID(id).Get(&existingNote); !has {
					fmt.Println(err)
					return cc.JSON(http.StatusConflict, map[string]interface{}{
						"Error":  "The snippet does not exist",
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
			}
		}

		return cc.JSON(http.StatusOK, map[string]interface{}{
			"Status": http.StatusOK,
		})
	}

	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"Error":  "Snippet not found",
		"Status": http.StatusNotFound,
	})
}
