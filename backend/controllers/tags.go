package controllers

import (
	"code-snippets/repository"
	. "code-snippets/types"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

const UNTAGGED = "[untagged]"

func GetTags(c echo.Context) error {
	tags := make([]NoteTag, 0)
	tags = append(tags, NoteTag{Alias: "[untagged]", ID: 0})

	if cc, ok := c.(*repository.CustomContext); ok {
		if db, err := cc.DB(); err == nil {
			_ = db.Find(&tags)

			for i := 1; i < len(tags); i++ {
				count, _ := db.Where("tag_id=?", tags[i].ID).Count(&Note{})
				tags[i].Notes = uint(count)
			}
		}
	}

	return c.JSON(http.StatusOK, tags)
}

func PostTag(c echo.Context) error {
	var (
		count int64
	)

	tag := new(NoteTag)

	if err := c.Bind(tag); err != nil {
		return c.JSON(http.StatusConflict, map[string]interface{}{
			"Error":  "Wrong parameters",
			"Status": http.StatusBadRequest,
		})
	}

	if strings.Trim(tag.Alias, " ") == "" {
		return c.JSON(http.StatusConflict, map[string]interface{}{
			"Error":  "Empty alias is not allowed",
			"Status": http.StatusBadRequest,
		})
	}

	if cc, ok := c.(*repository.CustomContext); ok {
		if db, err := cc.DB(); err == nil {
			if tag.ID <= 0 {
				// Create a new tag
				count, _ = db.Where("alias=?", tag.Alias).Count(&NoteTag{})
				if count > 0 || tag.Alias == UNTAGGED {
					return cc.JSON(http.StatusConflict, map[string]interface{}{
						"Error":  "Tag name already exists",
						"Status": http.StatusConflict,
					})
				}

				if _, err = db.InsertOne(*tag); err != nil {
					return cc.JSON(http.StatusConflict, map[string]interface{}{
						"Error":  err.Error(),
						"Status": http.StatusNotFound,
					})
				}
			} else {
				// Update an existing tag
				if count, _ = db.Where("alias=? AND id<>?", tag.Alias, tag.ID).Count(&NoteTag{}); count > 0 {
					return cc.JSON(http.StatusConflict, map[string]interface{}{
						"Error":  "Tag name already taken",
						"Status": http.StatusConflict,
					})
				}

				if _, err = db.ID(tag.ID).Update(*tag); err != nil {
					return cc.JSON(http.StatusConflict, map[string]interface{}{
						"Error":  err.Error(),
						"Status": http.StatusNotFound,
					})
				}
			}
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"Status": http.StatusOK,
	})
}

func GetTag(c echo.Context) error {
	if id, err := strconv.ParseInt(c.Param("id"), BASE10, strconv.IntSize); err == nil {
		var tag NoteTag

		if cc, ok := c.(*repository.CustomContext); ok {
			if db, err := cc.DB(); err == nil {
				if ok, err = db.ID(id).Get(&tag); ok && err == nil {
					return c.JSON(http.StatusOK, tag)
				}

				return c.JSON(http.StatusNotFound, map[string]interface{}{
					"Error":  err.Error(),
					"Status": http.StatusNotFound,
				})
			}
		}
	}

	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"Error":  "Tag not found",
		"Status": http.StatusNotFound,
	})
}

func DeleteTag(c echo.Context) error {
	var (
		count int64
	)

	if id, err := strconv.ParseInt(c.Param("id"), BASE10, strconv.IntSize); err == nil {
		if cc, ok := c.(*repository.CustomContext); ok {
			if db, err := cc.DB(); err == nil {
				// Check if there are notes for the tag
				if count, _ = db.Where("tag_id=?", id).Count(&Note{}); count > 0 {
					return cc.JSON(http.StatusConflict, map[string]interface{}{
						"Error":  "Tag has notes. Remove them first.",
						"Status": http.StatusConflict,
					})
				}

				if _, err = db.ID(id).Delete(&NoteTag{}); err == nil {
					return c.JSON(http.StatusOK, map[string]interface{}{
						"Status": http.StatusOK,
					})
				} else {
					return c.JSON(http.StatusNotFound, map[string]interface{}{
						"Error":  err.Error(),
						"Status": http.StatusBadRequest,
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
