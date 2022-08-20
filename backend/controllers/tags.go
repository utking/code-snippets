package controllers

import (
	"code-snippets/repository"
	. "code-snippets/types"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

const untagged = "[untagged]"

func GetTags(c echo.Context) error {
	tags := make([]NoteTag, 0)
	tags = append(tags, NoteTag{Alias: untagged, Snippets: 0})
	uniqTags := make([]Note, 0)
	userID, _ := GetUserID(c)

	if cc, ok := c.(*repository.CustomContext); ok {
		if db, err := cc.DB(); err == nil {
			if err = db.Where("user_id=?", userID).Distinct("tag").Find(&uniqTags); err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"Error":  err.Error(),
					"Status": http.StatusInternalServerError,
				})
			}

			for _, tag := range uniqTags {
				count, _ := db.Where("tag=?", tag.Tag).Count(&Note{})

				tags = append(tags, NoteTag{Alias: tag.Tag, Snippets: uint64(count)})
			}
		}
	}

	return c.JSON(http.StatusOK, tags)
}

// Update tag
func PutTag(c echo.Context) error {
	var (
		count int64
	)

	curTag := strings.Trim(c.Param("tag"), " ")
	userID, _ := GetUserID(c)

	tag := new(NoteTag)
	note := new(Note)

	if err := c.Bind(tag); err != nil {
		return c.JSON(http.StatusConflict, map[string]interface{}{
			"Error":  "Wrong parameters",
			"Status": http.StatusBadRequest,
		})
	}

	note.Tag = strings.Trim(tag.Alias, " ")
	note.UserID = userID

	if tag.Alias == "" {
		return c.JSON(http.StatusConflict, map[string]interface{}{
			"Error":  "Empty alias is not allowed",
			"Status": http.StatusBadRequest,
		})
	}

	if tag.Alias == curTag {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"Status": http.StatusOK,
		})
	}

	if cc, ok := c.(*repository.CustomContext); ok {
		if db, err := cc.DB(); err == nil {
			// Update an existing tag
			// Find there is an existing tag first
			if count, _ = db.Where("tag=? AND user_id=?", tag.Alias, userID).Count(&NoteTag{}); count > 0 {
				return cc.JSON(http.StatusConflict, map[string]interface{}{
					"Error":  "Tag name already taken",
					"Status": http.StatusConflict,
				})
			}

			if _, err = db.Where("tag=? AND user_id=?", curTag, userID).Update(*note); err != nil {
				return cc.JSON(http.StatusConflict, map[string]interface{}{
					"Error":  err.Error(),
					"Status": http.StatusNotFound,
				})
			}
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"Status": http.StatusOK,
	})
}
