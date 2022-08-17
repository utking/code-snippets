package controllers

import (
	"code-snippets/repository"
	. "code-snippets/types"
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"xorm.io/xorm"
)

func GetNote(c echo.Context) error {
	if id, err := strconv.ParseInt(c.Param("id"), BASE10, strconv.IntSize); err == nil {
		var (
			note Note
		)

		if cc, ok := c.(*repository.CustomContext); ok {
			if db, err := cc.DB(); err == nil {
				if ok, err = db.ID(id).Get(&note); err == nil && ok {
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

func GetTagNotes(c echo.Context) error {
	notes := make(Notes, 0)

	if tag := strings.Trim(c.Param("tag"), " "); tag != "" {
		if cc, ok := c.(*repository.CustomContext); ok {
			if db, err := cc.DB(); err == nil {
				if err = db.Find(&notes, &Note{Tag: tag}); err != nil {
					return c.JSON(http.StatusOK, map[string]interface{}{
						"Error":  err.Error(),
						"Status": http.StatusNotFound,
					})
				}
			}
		}
	}

	return c.JSON(http.StatusOK, notes)
}

func DeleteNote(c echo.Context) error {
	if id, err := strconv.ParseInt(c.Param("id"), BASE10, strconv.IntSize); err == nil {
		if cc, ok := c.(*repository.CustomContext); ok {
			if db, err := cc.DB(); err == nil {
				if _, err = db.ID(id).Delete(&Note{}); err == nil {
					return c.JSON(http.StatusOK, map[string]interface{}{
						"Status": http.StatusOK,
					})
				}

				return c.JSON(http.StatusNotFound, map[string]interface{}{
					"Error":  err.Error(),
					"Status": http.StatusNotFound,
				})
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

	if err := c.Bind(note); err != nil {
		return c.JSON(http.StatusConflict, map[string]interface{}{
			"Error":  "Wrong parameters",
			"Status": http.StatusBadRequest,
		})
	}

	if strings.Trim(note.Tag, " ") == "" {
		return c.JSON(http.StatusConflict, map[string]interface{}{
			"Error":  "Wrong or missing tag",
			"Status": http.StatusBadRequest,
		})
	}

	if strings.Trim(note.Title, " ") == "" || strings.Trim(note.Content, " ") == "" {
		return c.JSON(http.StatusConflict, map[string]interface{}{
			"Error":  "Empty title or content is not allowed",
			"Status": http.StatusBadRequest,
		})
	}

	if cc, ok := c.(*repository.CustomContext); ok {
		if db, err := cc.DB(); err == nil {
			if count, _ = db.Where("tag=? AND title=?", note.Tag, note.Title).Count(&Note{}); count > 0 {
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

	return c.JSON(http.StatusConflict, map[string]interface{}{
		"Error":  "Bad context",
		"Status": http.StatusNotFound,
	})
}

func PutNote(c echo.Context) error {
	var (
		existingNote Note
		db           *xorm.Engine
	)

	if id, err := strconv.ParseInt(c.Param("id"), BASE10, strconv.IntSize); err == nil {
		newNote := new(Note)

		if cc, ok := c.(*repository.CustomContext); ok {
			if err = cc.Bind(newNote); err != nil {
				return cc.JSON(http.StatusConflict, map[string]interface{}{
					"Error":  "Wrong parameters",
					"Status": http.StatusBadRequest,
				})
			}

			newNote.ID = uint16(id)

			if strings.Trim(newNote.Tag, " ") == " " {
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
				if _, err = db.ID(id).Get(&existingNote); err != nil {
					return cc.JSON(http.StatusConflict, map[string]interface{}{
						"Error":  err.Error(),
						"Status": http.StatusNotFound,
					})
				}

				if _, err = db.ID(id).Update(*newNote); err != nil {
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
