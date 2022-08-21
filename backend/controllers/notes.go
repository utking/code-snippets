package controllers

import (
	"code-snippets/repository"
	. "code-snippets/types"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"xorm.io/xorm"
)

func GetNote(c echo.Context) error {
	if id, err := strconv.ParseInt(c.Param("id"), BASE10, strconv.IntSize); err == nil {
		var (
			note    Note
			exists  bool
			toShare SharedNote
		)

		userID, _ := GetUserID(c)

		if cc, ok := c.(*repository.CustomContext); ok {
			if db, err := cc.DB(); err == nil {
				if ok, err = db.ID(id).Where("user_id=?", userID).Get(&note); err == nil && ok {
					if exists, _ = db.Where("note_id=? AND user_id=?", id, userID).Get(&toShare); exists {
						note.ShareHash = toShare.Hash
					}

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

func ShareNote(c echo.Context) error {
	var (
		note    Note
		toShare SharedNote
		db      *xorm.Engine
		err     error
	)

	userID, _ := GetUserID(c)
	shareNote := new(SharedNote)

	if err = c.Bind(shareNote); err != nil {
		return c.JSON(http.StatusConflict, map[string]interface{}{
			"Error":  "Wrong parameters",
			"Status": http.StatusBadRequest,
		})
	}

	shareNote.UserID = userID
	shareNote.Hash = shareNote.CalcHash()

	err = fmt.Errorf("snippet not found")

	if cc, exists := c.(*repository.CustomContext); exists {
		if db, err = cc.DB(); err == nil {
			// Check the snippet exists
			if exists, err = db.Where("id=? AND user_id=?", shareNote.NoteID, shareNote.UserID).Get(&note); !exists {
				return c.JSON(http.StatusNotFound, map[string]interface{}{
					"Error":  "The snippet does not exist",
					"Status": http.StatusBadRequest,
				})
			} else if err != nil {
				return c.JSON(http.StatusNotFound, map[string]interface{}{
					"Error":  err.Error(),
					"Status": http.StatusBadRequest,
				})
			}

			// Check if it is already shared
			if exists, err = db.Where("note_id=? AND user_id=?", shareNote.NoteID, shareNote.UserID).Get(&toShare); exists {
				note.ShareHash = toShare.Hash

				return c.JSON(http.StatusOK, note)
			} else if err != nil {
				return c.JSON(http.StatusNotFound, map[string]interface{}{
					"Error":  err.Error(),
					"Status": http.StatusBadRequest,
				})
			}

			// All is good, try to share it now
			if _, err = db.InsertOne(*shareNote); err != nil {
				return cc.JSON(http.StatusConflict, map[string]interface{}{
					"Error":  err.Error(),
					"Status": http.StatusNotFound,
				})
			}

			return cc.JSON(http.StatusOK, note)
		}
	}

	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"Error":  err.Error(),
		"Status": http.StatusNotFound,
	})
}

func GetSharedNote(c echo.Context) error {
	hash := strings.Trim(c.Param("hash"), " ")

	if hash != "" {
		var (
			sharedNote SharedNote
			note       Note
		)

		sharedNote.Hash = hash

		if cc, ok := c.(*repository.CustomContext); ok {
			if db, err := cc.DB(); err == nil {
				if ok, err = db.Where("hash=?", hash).Get(&sharedNote); err == nil && ok {
					if ok, err = db.Where("id=?", sharedNote.NoteID).Get(&note); err == nil && ok {
						return c.JSON(http.StatusOK, note)
					}
				}
			}
		}
	}

	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"Error":  "Snippet was not found",
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

				// Delete its share, if exists
				_, _ = db.Where("user_id=?", userID).Delete(&SharedNote{NoteID: uint16(id)})

				// Now delete the note
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

// Create a new snippet
func PostNote(c echo.Context) error {
	var (
		count, snippetID int64
		result           sql.Result
	)

	userID, _ := GetUserID(c)
	snippet := new(Note)

	if err := c.Bind(snippet); err != nil {
		return c.JSON(http.StatusConflict, map[string]interface{}{
			"Error":  "Wrong parameters",
			"Status": http.StatusBadRequest,
		})
	}

	snippet.Tag = strings.Trim(snippet.Tag, " ")
	snippet.Title = strings.Trim(snippet.Title, " ")
	snippet.UserID = userID

	if snippet.Tag == "" || snippet.Tag == untagged {
		return c.JSON(http.StatusConflict, map[string]interface{}{
			"Error":  "Wrong or missing tag",
			"Status": http.StatusBadRequest,
		})
	}

	if strings.Trim(snippet.Title, " ") == "" || strings.Trim(snippet.Content, " ") == "" {
		return c.JSON(http.StatusConflict, map[string]interface{}{
			"Error":  "Empty title or content is not allowed",
			"Status": http.StatusBadRequest,
		})
	}

	if cc, ok := c.(*repository.CustomContext); ok {
		if db, err := cc.DB(); err == nil {
			if count, _ = db.Count(*snippet); count > 0 {
				return cc.JSON(http.StatusConflict, map[string]interface{}{
					"Error":  "The snippet with this alias exists for the tag",
					"Status": http.StatusConflict,
				})
			}

			if _, err = db.InsertOne(*snippet); err != nil {
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

			if snippetID, err = result.LastInsertId(); err != nil {
				return cc.JSON(http.StatusConflict, map[string]interface{}{
					"Error":  err.Error(),
					"Status": http.StatusNotFound,
				})
			}

			var newSnippet Note

			if _, err = db.ID(snippetID).Get(&newSnippet); err == nil {
				return c.JSON(http.StatusOK, newSnippet)
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
		existing Note
		db       *xorm.Engine
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
				if _, err = db.ID(id).Where("user_id=?", userID).Get(&existing); err != nil {
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

				var newSnippet Note

				if _, err = db.ID(id).Get(&newSnippet); err == nil {
					return c.JSON(http.StatusOK, newSnippet)
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
