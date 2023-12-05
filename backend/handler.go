package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
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

// POST /items/sync replaces the complete list
func syncItems(db *sql.DB, notifier *Notifier) echo.HandlerFunc {
	return func(ctx echo.Context) error {

		// bind body into struct
		list := &ItemCollection{}
		err := ctx.Bind(&list)
		if err != nil {
			ctx.Logger().Infof("replaceItemList: Bind Error with request %v: %v", ctx.Request().Body, err)
			return echo.NewHTTPError(http.StatusBadRequest, "Wrong Input")
		}
		// Compare version from input with version from DB
		versions, err := GetVersions(db)
		if err != nil {
			ctx.Logger().Infof("replaceItemList: Cannot get versions from DB %v", err)
			return echo.NewHTTPError(http.StatusBadRequest, "DB Error")
		}

		// if version from API is newer as version from DB => replace
		if list.Version > versions.ItemVersion {

			// do database operation
			err = ReplaceItemList(db, list)
			if err != nil {
				ctx.Logger().Infof("syncItems: Database Error on replace %v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Could not change item")
			}

			// looks fine, notify all the listening clients:
			notifier.Send("UPDATE")

		}

		items, err := GetAllItems(db)
		if err != nil {
			ctx.Logger().Infof("syncItem: Database Error on get %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not change item")
		}
		return ctx.JSON(http.StatusOK, items)
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

		go func() {
			<-ctx.Request().Context().Done()
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

// POST /shops/sync replaces the complete list
func syncShops(db *sql.DB, notifier *Notifier) echo.HandlerFunc {
	return func(ctx echo.Context) error {

		// bind body into struct
		list := &ShopCollection{}
		err := ctx.Bind(&list)
		if err != nil {
			ctx.Logger().Infof("replaceItemList: Bind Error with request %v: %v", ctx.Request().Body, err)
			return echo.NewHTTPError(http.StatusBadRequest, "Wrong Input")
		}

		// Compare version from input with version from DB
		versions, err := GetVersions(db)
		if err != nil {
			ctx.Logger().Infof("replaceItemList: Cannot get versions from DB %v", err)
			return echo.NewHTTPError(http.StatusBadRequest, "DB Error")
		}

		// if version from API is newer as version from DB => replace
		if list.Version > versions.ShopVersion {

			// do database operation
			err = ReplaceShopList(db, list)
			if err != nil {
				ctx.Logger().Infof("syncShop: Database Error on replace %v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "DB Error")

			}

			// looks fine, notify all the listening clients:
			notifier.Send("UPDATE")
		}

		shops, err := GetAllShops(db)
		if err != nil {
			ctx.Logger().Infof("syncShop: Database Error on get %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "DB Error")
		}
		return ctx.JSON(http.StatusOK, shops)

	}
}
