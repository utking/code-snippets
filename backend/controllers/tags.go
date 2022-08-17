package controllers

import (
	"code-snippets/repository"
	. "code-snippets/types"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

const untagged = "[untagged]"

func GetTags(c echo.Context) error {
	tags := make([]NoteTag, 0)
	tags = append(tags, NoteTag{Alias: untagged, Snippets: 0})
	uniqTags := make([]Note, 0)

	if cc, ok := c.(*repository.CustomContext); ok {
		if db, err := cc.DB(); err == nil {
			if err = db.Distinct("tag").Find(&uniqTags); err != nil {
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
