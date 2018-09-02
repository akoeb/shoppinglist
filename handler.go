package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

// Handler

// GET / shows list of all items
func showAllItems(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		items, err := GetAllItems(db)
		if err != nil {
			ctx.Logger().Infof("showAllItems: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not read items")
		}
		return ctx.JSON(http.StatusOK, items)
	}
}

// GET /:id/ shows one item
func showOneItem(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter")
		}
		item, err := GetItemByID(db, id)
		if err != nil {
			ctx.Logger().Infof("showOneItem: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not read item")
		}
		return ctx.JSON(http.StatusOK, item)
	}
}

// POST / creates a new item
func createItem(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		item := &Item{}

		// bind body into struct
		err := ctx.Bind(item)
		if err != nil {
			ctx.Logger().Infof("createItem: ind Error with request %v: %v", ctx.Request().Body, err)
			return echo.NewHTTPError(http.StatusBadRequest, "Wrong Input")
		}

		// field validation
		if ok, errors := item.Valid(); !ok {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Wrong Input: %v", errors))
		}

		// creates can not have an id
		if item.ID > 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "Can not create item that has an id")
		}

		// write to db
		err = UpsertItem(db, item)
		if err != nil {
			ctx.Logger().Infof("createItem: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not write item")
		}

		// return the modified item with new id:
		return ctx.JSON(http.StatusOK, item)
	}
}

// PUT /:id/ updates an existing item
func updateItem(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter")
		}
		item := &Item{}
		err = ctx.Bind(&item)
		if err != nil {
			ctx.Logger().Infof("updateItem: Bind Error with request %v: %v", ctx.Request().Body, err)
			return echo.NewHTTPError(http.StatusBadRequest, "Wrong Input")
		}
		if item.ID != id {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("item with id %d can not be updated in path with id %d", item.ID, id))
		}
		err = UpsertItem(db, item)
		if err != nil {
			ctx.Logger().Infof("updateItem: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not change item")

		}
		return ctx.JSON(http.StatusOK, item)
	}
}

// DELETE /:id/ deletes one item
func deleteOneItem(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter")
		}
		_, err = DeleteItemByID(db, id)
		if err != nil {
			ctx.Logger().Infof("deleteItem: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not delete item")

		}
		return ctx.NoContent(http.StatusNoContent)
	}
}

func reorderItems(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// this request has only a map[int]int as payload:
		var objmap map[int]int
		err := ctx.Bind(&objmap)
		if err != nil {
			ctx.Logger().Infof("reorderItems: Can not deserialize to map %v", err)
			return echo.NewHTTPError(http.StatusBadRequest, "wrong payload")

		}
		err = reOrderItems(db, objmap)
		if err != nil {
			ctx.Logger().Infof("reorderItems: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not reorder items")

		}
		items, err := GetAllItems(db)
		if err != nil {
			ctx.Logger().Infof("reorderItems: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not read items")
		}
		return ctx.JSON(http.StatusOK, items)
	}
}

// not yet implemented
func deleteManyItems(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "NOT YET IMPLEMENTED")
	}
}
