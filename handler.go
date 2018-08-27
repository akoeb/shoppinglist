package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

// Handler
func showAllItems(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		items, err := GetAllItems(db)
		if err != nil {
			return err
		}
		return ctx.JSON(http.StatusOK, items)
	}
}

func showOneItem(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return err
		}
		item, err := GetItemById(db, id)
		if err != nil {
			return err
		}
		return ctx.JSON(http.StatusOK, item)
	}
}

func createItem(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		item := &Item{}
		err := ctx.Bind(item)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Wrong Input")
		}
		if item.ID > 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "Can not create item with fixed id")
		}
		err = UpsertItem(db, item)
		if err != nil {
			return err
		}
		return ctx.JSON(http.StatusOK, item)
	}
}

func updateItem(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return err
		}
		item := &Item{}
		err = ctx.Bind(&item)
		if err != nil {
			return err
		}
		if item.ID != id {
			return fmt.Errorf("item with id %d can not be updated in path with id %d", item.ID, id)
		}
		err = UpsertItem(db, item)
		if err != nil {
			return err
		}
		return ctx.JSON(http.StatusOK, item)
	}
}
func deleteOneItem(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return err
		}
		_, err = DeleteItemById(db, id)
		if err != nil {
			return err
		}
		return ctx.NoContent(http.StatusNoContent)
	}
}

func deleteManyItems(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "NOT YET IMPLEMENTED")
	}
}
