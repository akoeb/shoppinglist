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

// POST /items/:id/shop/:sid assigns item :id to shop :sid
func assignItemToShop(db *sql.DB, notifier *Notifier) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// get parameter
		itemId, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter itemid")
		}

		shopId, err := strconv.Atoi(ctx.Param("sid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter shopid")
		}
		item, err := GetItemByID(db, itemId)
		if err != nil {
			ctx.Logger().Infof("assignItemtoShop: Database Error in getting Item %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not read item")
		}
		shop, err := GetShopByID(db, shopId)
		if err != nil {
			ctx.Logger().Infof("assignItemtoShop: Database Error in getting Shop %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not read shop")
		}

		// overwrite shop in item:
		item.Shop = &shop

		// save to database
		err = UpsertItem(db, &item)
		if err != nil {
			ctx.Logger().Infof("updateItem: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not change item")

		}

		// looks fine, notify all the listening clients:
		notifier.Send("UPDATE")

		// and inform the client
		return ctx.NoContent(http.StatusNoContent)
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

// GET /shops - show all shops registered
func showAllShops(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		shop, err := GetAllShops(db)
		if err != nil {
			ctx.Logger().Infof("showAllShops: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not read shops")
		}
		return ctx.JSON(http.StatusOK, shop)
	}
}

// GET /shops/:id - get one shop
func showOneShop(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter")
		}
		shop, err := GetShopByID(db, id)
		if err != nil {
			ctx.Logger().Infof("showOneShop: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not read shop")
		}
		return ctx.JSON(http.StatusOK, shop)
	}
}

// POST /shops create a shop
func createShop(db *sql.DB, notifier *Notifier) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		shop := &Shop{}

		// bind body into struct
		err := ctx.Bind(shop)
		if err != nil {
			ctx.Logger().Infof("createShop: Bind Error with request %v: %v", ctx.Request().Body, err)
			return echo.NewHTTPError(http.StatusBadRequest, "Wrong Input")
		}

		// creates can not have an id
		if shop.ID > 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "Can not create shop that has an id")
		}

		// write to db
		err = UpsertShop(db, shop)
		if err != nil {
			ctx.Logger().Infof("createShop: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not write shop")
		}

		// looks fine, notify all the listening clients:
		notifier.Send("UPDATE")

		// return the modified item with new id:
		return ctx.JSON(http.StatusOK, shop)
	}
}

// PUT /shops/:id modify a shop
func updateShop(db *sql.DB, notifier *Notifier) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// read parameter
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter")
		}

		// bind body into struct
		shop := &Shop{}
		err = ctx.Bind(&shop)
		if err != nil {
			ctx.Logger().Infof("updateShop: Bind Error with request %v: %v", ctx.Request().Body, err)
			return echo.NewHTTPError(http.StatusBadRequest, "Wrong Input")
		}
		// some validation
		if shop.ID != id {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("shop with id %d can not be updated in path with id %d", shop.ID, id))
		}

		// do database operation
		err = UpsertShop(db, shop)
		if err != nil {
			ctx.Logger().Infof("updateItem: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not change item")

		}

		// looks fine, notify all the listening clients:
		notifier.Send("UPDATE")

		// inform the client
		return ctx.JSON(http.StatusOK, shop)
	}
}

// DELETE /shops/:id - delete a shop
func deleteOneShop(db *sql.DB, notifier *Notifier) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// get parameter
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter")
		}

		// do the database magic
		_, err = DeleteShopByID(db, id)
		if err != nil {
			ctx.Logger().Infof("deleteShop: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not delete shop")
		}

		// looks fine, notify all the listening clients:
		notifier.Send("UPDATE")

		// and inform the client
		return ctx.NoContent(http.StatusNoContent)
	}
}
