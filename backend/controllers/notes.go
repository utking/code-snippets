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

		userID, _ := GetUserID(c)

		if cc, ok := c.(*repository.CustomContext); ok {
			if db, err := cc.DB(); err == nil {
				if ok, err = db.ID(id).Where("user_id=?", userID).Get(&note); err == nil && ok {
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
	userID, _ := GetUserID(c)

	if tag := strings.Trim(c.Param("tag"), " "); tag != "" {
		if cc, ok := c.(*repository.CustomContext); ok {
			if db, err := cc.DB(); err == nil {
				if err = db.Find(&notes, &Note{Tag: tag, UserID: userID}); err != nil {
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
				userID, _ := GetUserID(c)

				if _, err = db.ID(id).Where("user_id=?", userID).Delete(&Note{}); err == nil {
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

// Create a new note
func PostNote(c echo.Context) error {
	var (
		count, noteID int64
		result        sql.Result
	)

	userID, _ := GetUserID(c)
	note := new(Note)

	if err := c.Bind(note); err != nil {
		return c.JSON(http.StatusConflict, map[string]interface{}{
			"Error":  "Wrong parameters",
			"Status": http.StatusBadRequest,
		})
	}

	note.Tag = strings.Trim(note.Tag, " ")
	note.Title = strings.Trim(note.Title, " ")
	note.UserID = userID

	if note.Tag == "" || note.Tag == untagged {
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
			if count, _ = db.Count(*note); count > 0 {
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
		userID, _ := GetUserID(c)

		if cc, ok := c.(*repository.CustomContext); ok {
			if err = cc.Bind(newNote); err != nil {
				return cc.JSON(http.StatusConflict, map[string]interface{}{
					"Error":  "Wrong parameters",
					"Status": http.StatusBadRequest,
				})
			}

			newNote.ID = uint16(id)
			newNote.UserID = userID

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
				if _, err = db.ID(id).Where("user_id=?", userID).Get(&existingNote); err != nil {
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
