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
func createItem(db *sql.DB, notifier *Notifier) echo.HandlerFunc {
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

		// looks fine, notify all the listening clients:
		notifier.Send("UPDATE")

		// return the modified item with new id:
		return ctx.JSON(http.StatusOK, item)
	}
}

// PUT /:id/ updates an existing item
func updateItem(db *sql.DB, notifier *Notifier) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// read parameter
		id, err := strconv.Atoi(ctx.Param("id"))
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
		if item.ID != id {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("item with id %d can not be updated in path with id %d", item.ID, id))
		}

		// do database operation
		err = UpsertItem(db, item)
		if err != nil {
			ctx.Logger().Infof("updateItem: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not change item")

		}

		// looks fine, notify all the listening clients:
		notifier.Send("UPDATE")

		// inform the client
		return ctx.JSON(http.StatusOK, item)
	}
}

// DELETE /:id/ deletes one item
func deleteOneItem(db *sql.DB, notifier *Notifier) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// get parameter
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter")
		}

		// do the database magic
		_, err = DeleteItemByID(db, id)
		if err != nil {
			ctx.Logger().Infof("deleteItem: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not delete item")
		}

		// looks fine, notify all the listening clients:
		notifier.Send("UPDATE")

		// and inform the client
		return ctx.NoContent(http.StatusNoContent)
	}
}

func reorderItems(db *sql.DB, notifier *Notifier) echo.HandlerFunc {
	return func(ctx echo.Context) error {

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
		err = reOrderItems(db, intmap)
		if err != nil {
			ctx.Logger().Infof("reorderItems: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not reorder items")

		}

		// looks fine, notify all the listening clients:
		notifier.Send("UPDATE")

		// and inform the client
		// read the actual list for sending back to the client
		items, err := GetAllItems(db)
		if err != nil {
			ctx.Logger().Infof("reorderItems: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not read items")
		}
		return ctx.JSON(http.StatusOK, items)
	}
}

// not yet implemented
func deleteManyItems(db *sql.DB, notifier *Notifier) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		requestedStatus := ctx.QueryParam("status")
		if !isAllowedStatusCode(requestedStatus) {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("wrong status code: %s", requestedStatus))
		}

		// do the database magic
		_, err := deleteManyItemsByStatus(db, requestedStatus)
		if err != nil {
			ctx.Logger().Infof("deleteManyItems: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not delete items")
		}

		// looks fine, notify all the listening clients:
		notifier.Send("UPDATE")

		// and inform the client
		return ctx.NoContent(http.StatusNoContent)
	}
}

// handle event streams:
func eventsStream(notifier *Notifier) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		ctx.Logger().Infof("eventsStream called")
		receiverID, err := notifier.NewReceiver()
		if err != nil {
			ctx.Logger().Errorf("eventsStream Error creating receiver: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)

		}

		// close notifier to clean up notifier channel if connection closes:

		closeNotify := ctx.Response().CloseNotify()
		go func() {
			<-closeNotify
			// connection close, do cleanup
			ctx.Logger().Infof("eventsStream: close notify triggered")
			notifier.RemoveReceiver(receiverID)
		}()

		ctx.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
		ctx.Response().WriteHeader(http.StatusOK)
		ctx.Response().Flush()
		for {
			// listen on channel:
			cmd := notifier.Listen(receiverID)

			// write message to client
			_, err := ctx.Response().Write([]byte(fmt.Sprintf("data: {\"cmd\": \"%s\"}\n\n", cmd)))
			if err != nil {
				notifier.RemoveReceiver(receiverID)
				msg := fmt.Sprintf("Error writing to stream: %v", err)
				ctx.Logger().Infof("eventsStream: Error %v", msg)
				return echo.NewHTTPError(http.StatusInternalServerError, msg)
			}
			ctx.Response().Flush()
		}
	}
}
