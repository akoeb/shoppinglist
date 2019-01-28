package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

// Handler for location endpoints inside category:
// apis.GET("/locations", showAllLocations(db))          // list all locations, can be filtered by category id
// apis.POST("/locations", createLocation(db, notifier)) // create a location
// apis.PUT("/locations/:locid", updateLocation(db, notifier))
// apis.DELETE("/locations/:locid", deleteOneLocation(db, notifier))
// apis.POST("/locations/reorder", reorderItemsByLocation(db, notifier))

// apis.GET("/locations", showAllLocations(db))
// list all locations, can be filtered by category id
func showAllLocations(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// do we have a category query string parameter for filtering?
		catid, err := strconv.Atoi(ctx.QueryParam("category"))
		if err != nil {
			// either the parameter is unset or empty, treat as no filter
			catid = 0
		}
		locations, err := GetLocations(db, catid)
		if err != nil {
			ctx.Logger().Infof("showAllLocationsInCategory: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not read location list")
		}
		return ctx.JSON(http.StatusOK, locations)
	}
}

// apis.POST("/locations", createLocation(db, notifier))
func createLocation(db *sql.DB, notifier *Notifier) echo.HandlerFunc {
	return func(ctx echo.Context) error {

		// bind body into struct
		location := &Location{}
		err := ctx.Bind(location)
		if err != nil {
			ctx.Logger().Infof("createLocation: bind Error with request %v: %v", ctx.Request().Body, err)
			return echo.NewHTTPError(http.StatusBadRequest, "Wrong Input")
		}

		// field validation
		if ok, errors := location.Valid(); !ok {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Wrong Input: %v", errors))
		}

		// creates can not have an id
		if location.ID > 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "Can not create location that has an id")
		}

		// write to db
		err = UpsertLocation(db, location)
		if err != nil {
			ctx.Logger().Infof("createLocation: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not write location")
		}

		// looks fine, notify all the listening clients:
		notifier.Send("UPDATE location", location.ID)

		// return the modified location with new id:
		return ctx.JSON(http.StatusOK, location)
	}
}

// apis.PUT("/locations/:locid", updateLocation(db, notifier))
func updateLocation(db *sql.DB, notifier *Notifier) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// read parameter
		id, err := strconv.Atoi(ctx.Param("locid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter")
		}

		// bind body into struct
		location := &Location{}
		err = ctx.Bind(&location)
		if err != nil {
			ctx.Logger().Infof("updateLocation: Bind Error with request %v: %v", ctx.Request().Body, err)
			return echo.NewHTTPError(http.StatusBadRequest, "Wrong Input")
		}
		// some validation
		if location.ID != id {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("location with id %d can not be updated in path with id %d", location.ID, id))
		}

		// do database operation
		err = UpsertLocation(db, location)
		if err != nil {
			ctx.Logger().Infof("updateLocation: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not change location")

		}

		// looks fine, notify all the listening clients:
		notifier.Send("UPDATE location", location.ID)

		// inform the client
		return ctx.JSON(http.StatusOK, location)
	}
}

// apis.DELETE("/locations/:locid", deleteOneLocation(db, notifier))
func deleteOneLocation(db *sql.DB, notifier *Notifier) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// get parameter
		id, err := strconv.Atoi(ctx.Param("locid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter")
		}

		// do the database magic
		_, err = DeleteLocationByID(db, id)
		if err != nil {
			ctx.Logger().Infof("deleteLocation: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not delete item")
		}

		// looks fine, notify all the listening clients:
		notifier.Send("UPDATE location", id)

		// and inform the client
		return ctx.NoContent(http.StatusNoContent)
	}
}

/*
TODO:
// apis.POST("/locations/:locid/reorder", reorderItemsByLocation(db, notifier))
func reorderItemsByLocation(db *sql.DB, notifier *Notifier) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// get parameter
		id, err := strconv.Atoi(ctx.Param("locid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter")
		}
		// this request has only a map[int]int as payload:
		// key is item id (as string, javascript facepalm), value is the new order no
		var strmap map[string]int
		err := ctx.Bind(&strmap)
		if err != nil {
			ctx.Logger().Infof("reorderItems: Can not deserialize to map %v", err)
			return echo.NewHTTPError(http.StatusBadRequest, "wrong payload")

		}

		// now we copy that map into one with int keys:
		var intmap map[int]int
		intmap = make(map[int]int)
		for k, v := range strmap {
			i, err := strconv.ParseInt(k, 10, 32)
			if err != nil {
				ctx.Logger().Infof("reorderItems: Can not convert to int-map %v", err)
				return echo.NewHTTPError(http.StatusBadRequest, "key is not a number: "+k)
			}
			intmap[int(i)] = v
		}

		// ok, do the update
		err = reorderItems(db, id, intmap)
		if err != nil {
			ctx.Logger().Infof("reorderItems: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not reorder items")

		}

		// looks fine, notify all the listening clients:
		notifier.Send("UPDATE location", id)

		// and inform the client
		// read the actual list for sending back to the client
		items, err := GetAllItemsByLocation(db)
		if err != nil {
			ctx.Logger().Infof("reorderItems: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not read items")
		}
		return ctx.JSON(http.StatusOK, items)

	}
}
*/
