package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Note struct {
	Title   string
	Syntax  string
	Content string
	TagID   uint16
	ID      uint16
	Indent  uint8
}

type Notes []Note

var notes Notes

const (
	indent = 4
)

func init() {
	notes = Notes{
		{ID: 1, Title: "Some important SELECT", Syntax: "mysql", TagID: 1,
			Content: "SELECT * from user", Indent: indent},
		{ID: 2, Title: "Drop ES index", Syntax: "shell", TagID: 2,
			Content: "curl -XDELETE localhost:9200/some-index", Indent: indent},
		{ID: 3, Title: "List indices", Syntax: "shell", TagID: 2,
			Content: "curl localhost:9200/_cat/indices?v", Indent: indent},
		{ID: 4, Title: "Some note", Syntax: "auto", TagID: 3,
			Content: "Just a note", Indent: indent},
	}
}

func GetNotes(c echo.Context) error {
	return c.JSON(http.StatusOK, notes)
}

func GetNote(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, strconv.IntSize)
	if err == nil {
		for _, itm := range notes {
			if itm.ID == uint16(id) {
				return c.JSON(http.StatusOK, itm)
			}
		}
	}

	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"Error":  "Note not found",
		"Status": http.StatusNotFound,
	})
}

func filterItems(items Notes, cmpFn func(Note) bool) Notes {
	result := make(Notes, 0)

	for _, item := range items {
		if cmpFn(item) {
			result = append(result, item)
		}
	}

	return result
}

func GetCategoryNotes(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, strconv.IntSize)

	if err == nil {
		return c.JSON(http.StatusOK, filterItems(notes,
			func(item Note) bool { return item.TagID == uint16(id) }))
	}

	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"Error":  "Note not found",
		"Status": http.StatusNotFound,
	})
}
