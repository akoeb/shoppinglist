package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

// Handler for categories for these endpoints
// apis.GET("/categories", showAllCategories(db))         // list all categories
// apis.GET("/categories/:catid", showOneCategory(db))    // show one category
// apis.POST("/categories", createCategory(db, notifier)) // create a category
// apis.PUT("/categories/:catid", updateCategory(db, notifier))
// apis.DELETE("/categories/:catid", deleteOneCategory(db, notifier))

// apis.GET("/categories", showAllCategories(db))
// list all categories
func showAllCategories(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		categories, err := GetAllCategories(db)
		if err != nil {
			ctx.Logger().Infof("showAllCategories: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not read categories")
		}
		return ctx.JSON(http.StatusOK, categories)
	}
}

// apis.GET("/categories/:catid", showOneCategory(db))
// show one category
func showOneCategory(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("catid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter")
		}
		category, err := GetCategoryByID(db, id)
		if err != nil {
			ctx.Logger().Infof("showOneCategory: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not read category")
		}
		return ctx.JSON(http.StatusOK, category)
	}
}

// apis.POST("/categories", createCategory(db, notifier)) // create a category
func createCategory(db *sql.DB, notifier *Notifier) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		category := &Category{}

		// bind body into struct
		err := ctx.Bind(category)
		if err != nil {
			ctx.Logger().Infof("createCategory: bind error with request %v: %v", ctx.Request().Body, err)
			return echo.NewHTTPError(http.StatusBadRequest, "Wrong Input")
		}

		// field validation
		if ok, errors := category.Valid(); !ok {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Wrong Input: %v", errors))
		}

		// creates can not have an id
		if category.ID > 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "Can not create category that has an id")
		}

		// write to db
		err = UpsertCategory(db, category)
		if err != nil {
			ctx.Logger().Infof("createCategory: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not write category")
		}

		// looks fine, notify all the listening clients:
		notifier.Send("UPDATE category", category.ID)

		// return the modified item with new id:
		return ctx.JSON(http.StatusOK, category)
	}
}

// apis.PUT("/categories/:catid", updateCategory(db, notifier))
func updateCategory(db *sql.DB, notifier *Notifier) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// read parameter
		id, err := strconv.Atoi(ctx.Param("catid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter")
		}

		// bind body into struct
		category := &Category{}
		err = ctx.Bind(&category)
		if err != nil {
			ctx.Logger().Infof("updateCategory: Bind Error with request %v: %v", ctx.Request().Body, err)
			return echo.NewHTTPError(http.StatusBadRequest, "Wrong Input")
		}
		// some validation
		if category.ID != id {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("category with id %d can not be updated in path with id %d", category.ID, id))
		}

		// do database operation
		err = UpsertCategory(db, category)
		if err != nil {
			ctx.Logger().Infof("updateCategory: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not change category")

		}

		// looks fine, notify all the listening clients:
		notifier.Send("UPDATE category", category.ID)

		// inform the client
		return ctx.JSON(http.StatusOK, category)
	}
}

// apis.DELETE("/categories/:catid", deleteOneCategory(db, notifier))
func deleteOneCategory(db *sql.DB, notifier *Notifier) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// get parameter
		catid, err := strconv.Atoi(ctx.Param("catid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Path Parameter")
		}

		// do the database magic
		_, err = DeleteCategoryByID(db, catid)
		if err != nil {
			ctx.Logger().Infof("deleteCategory: Database Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not delete category")
		}

		// looks fine, notify all the listening clients:
		notifier.Send("UPDATE category", catid)

		// and inform the client
		return ctx.NoContent(http.StatusNoContent)
	}
}
