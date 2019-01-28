package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

// Handler for items

// 	apis.GET("/categories/:catid/items", showAllItemsInCategory(db)) // list all items in one category, can be filtered by location id
//	apis.POST("/categories/:catid/items", createItem(db, notifier))  // create an item
//  apis.GET("/categories/:catid/items/:itemid", showOneItemInCategory(db, notifier))
//	apis.PUT("/categories/:catid/items/:itemid", updateItem(db, notifier))
//	apis.DELETE("/categories/:catid/items/:itemid", deleteOneItem(db, notifier))

// apis.GET("/categories/:catid/items", showAllItemsInCategory(db)) shows list of all items
// TODO: possibly filtered by location id
func showAllItemsInCategory(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		catid, err := strconv.Atoi(ctx.Param("catid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter Category ID")
		}
		// location filter active?
		locationID, err := strconv.Atoi(ctx.QueryParam("loc"))
		if err != nil {
			// no filter on error
			locationID = 0
		}

		items, err := GetItemsInCategory(db, catid, locationID)
		if err != nil {
			ctx.Logger().Infof("GetItemsInCategory: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not read items")
		}
		return ctx.JSON(http.StatusOK, items)
	}
}

// apis.GET("/categories/:catid/items/:itemid", showOneItemInCategory(db, notifier))
func showOneItemInCategory(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		catid, err := strconv.Atoi(ctx.Param("catid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter Category ID")
		}
		itemid, err := strconv.Atoi(ctx.Param("itemid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter Item ID")
		}
		item, err := GetItemByIDAndCategory(db, catid, itemid)
		if err != nil {
			ctx.Logger().Infof("showOneItem: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not read item")
		}
		if item.ID == 0 {
			return echo.NewHTTPError(http.StatusNotFound, "No item by that ID in this category")
		}
		return ctx.JSON(http.StatusOK, item)
	}
}

//	apis.POST("/categories/:catid/items", createItem(db, notifier))  // create an item
func createItem(db *sql.DB, notifier *Notifier) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		catid, err := strconv.Atoi(ctx.Param("catid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter Category ID")
		}
		if catid < 1 {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter Category ID")
		}

		item := &Item{}

		// bind body into struct
		err = ctx.Bind(item)
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
		err = UpsertItem(db, catid, item)
		if err != nil {
			ctx.Logger().Infof("createItem: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not write item")
		}

		// looks fine, notify all the listening clients:
		notifier.Send("UPDATE category", catid)

		// return the modified item with new id:
		return ctx.JSON(http.StatusOK, item)
	}
}

//	apis.PUT("/categories/:catid/items/:itemid", updateItem(db, notifier))
func updateItem(db *sql.DB, notifier *Notifier) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		catid, err := strconv.Atoi(ctx.Param("catid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter Category ID")
		}
		// read parameter
		itemid, err := strconv.Atoi(ctx.Param("itemid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter")
		}

		// bind body into struct
		item := &Item{}
		err = ctx.Bind(&item)
		if err != nil {
			ctx.Logger().Infof("updateItem: Bind Error with request %v: %v", ctx.Request().Body, err)
			return echo.NewHTTPError(http.StatusBadRequest, "Wrong Input")
		}
		// some validation
		if item.ID != itemid {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("item with id %d can not be updated in path with id %d", item.ID, itemid))
		}

		// do database operation
		err = UpsertItem(db, catid, item)
		if err != nil {
			ctx.Logger().Infof("updateItem: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not change item")

		}

		// looks fine, notify all the listening clients:
		notifier.Send("UPDATE category", catid)

		// inform the client
		return ctx.JSON(http.StatusOK, item)
	}
}

//	apis.DELETE("/categories/:catid/items/:itemid", deleteOneItem(db, notifier))
func deleteOneItem(db *sql.DB, notifier *Notifier) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// get parameter
		catid, err := strconv.Atoi(ctx.Param("catid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter Category ID")
		}
		itemid, err := strconv.Atoi(ctx.Param("itemid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter")
		}

		// do the database magic
		_, err = DeleteItemByID(db, catid, itemid)
		if err != nil {
			ctx.Logger().Infof("deleteItem: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not delete item")
		}

		// looks fine, notify all the listening clients:
		notifier.Send("UPDATE category", catid)

		// and inform the client
		return ctx.NoContent(http.StatusNoContent)
	}
}
