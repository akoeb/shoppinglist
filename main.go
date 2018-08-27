package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// Database
	db := initDB("storage.db")
	migrate(db)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	// Routes
	e.GET("/", showAllItems(db))
	e.GET("/:id", showOneItem(db))
	e.POST("/", createItem(db))
	e.PUT("/:id", updateItem(db))
	e.DELETE("/:id", deleteOneItem(db))
	e.DELETE("/", deleteManyItems(db))

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
