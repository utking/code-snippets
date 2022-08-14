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

const UNTAGGED = "[untagged]"

type NoteTag struct {
	Alias string `xorm:"varchar(32) not null"`
	ID    uint16 `xorm:"id pk autoincr"`
	Notes uint   `xorm:"-"`
}

func GetTags(c echo.Context) error {
	tags := make([]NoteTag, 0)
	tags = append(tags, NoteTag{Alias: "[untagged]", ID: 0})

	cc, ok := c.(*CustomContext)

	if ok {
		if db, err := cc.DB(); err == nil {
			_ = db.CreateTables(NoteTag{})
			_ = db.SQL("SELECT * from note_tag").Find(&tags)

			for i := 1; i < len(tags); i++ {
				count, _ := db.Where("tag_id=?", tags[i].ID).Count(&Note{})
				tags[i].Notes = uint(count)
			}
		}
	}

	return cc.JSON(http.StatusOK, tags)
}

func PostTag(c echo.Context) error {
	var (
		count int64
		// result     int64
		// existingID int16
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

			if tag.ID <= 0 {
				// Create a new tag
				count, _ = db.Where("alias=?", tag.Alias).Count(&NoteTag{})
				if count > 0 || tag.Alias == UNTAGGED {
					return cc.JSON(http.StatusConflict, map[string]interface{}{
						"Error":  "Tag name already exists",
						"Status": http.StatusConflict,
					})
				}

				_, err = db.InsertOne(*tag)

				if err != nil {
					return cc.JSON(http.StatusConflict, map[string]interface{}{
						"Error":  err.Error(),
						"Status": http.StatusNotFound,
					})
				}
			} else {
				// Update an existing tag
				count, _ = db.Where("alias=? AND id<>?", tag.Alias, tag.ID).Count(&NoteTag{})

				if count > 0 {
					return cc.JSON(http.StatusConflict, map[string]interface{}{
						"Error":  "Tag name already taken",
						"Status": http.StatusConflict,
					})
				}

				_, err = db.ID(tag.ID).Update(*tag)

				if err != nil {
					return cc.JSON(http.StatusConflict, map[string]interface{}{
						"Error":  err.Error(),
						"Status": http.StatusNotFound,
					})
				}
			}

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
	var (
		count int64
	)
	id, err := strconv.ParseInt(c.Param("id"), 10, strconv.IntSize)

	if err == nil {
		cc, ok := c.(*CustomContext)

		if ok {
			if db, err := cc.DB(); err == nil {
				_ = db.CreateTables(NoteTag{})

				// Check if there are notes for the tag
				count, _ = db.Where("tag_id=?", id).Count(&Note{})

				if count > 0 {
					return cc.JSON(http.StatusConflict, map[string]interface{}{
						"Error":  "Tag has notes. Remove them first.",
						"Status": http.StatusConflict,
					})
				}

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
