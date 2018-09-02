package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

// Options hold global application options that can be set via CLI
type Options struct {
	DatabaseFile         *string
	HTTPBaseAuthUser     *string
	HTTPBaseAuthPassword *string
	Port                 *int
}

func parseFlags() *Options {

	options := &Options{}

	// options:
	options.Port = flag.Int("port", 8080, "The Port that the application uses for listening")
	options.DatabaseFile = flag.String("db", "storage.db", "The file to store the sqlite3 database")
	options.HTTPBaseAuthUser = flag.String("user", "", "The user name for HTTP Base authentication")
	options.HTTPBaseAuthPassword = flag.String("password", "", "The password for HTTP Base authentication")

	// parse command line into options
	flag.Parse()

	// CLI Validation goes here
	if *options.HTTPBaseAuthPassword == "" && *options.HTTPBaseAuthUser != "" {
		log.Fatal("Can not use HTTP Base Authentication with only user, needs also password")
	}
	if *options.HTTPBaseAuthPassword != "" && *options.HTTPBaseAuthUser == "" {
		log.Fatal("Can not use HTTP Base Authentication with only password, needs also user")
	}

	return options
}

func main() {

	options := parseFlags()

	// Echo instance
	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	// Database
	db := initDB(*options.DatabaseFile)
	migrate(db)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS header
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	// routes for static files
	e.Static("/", "public")

	// apis have their own middlewares: group them
	apis := e.Group("/items")

	// only allow application/json content type:
	apis.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if ctx.Request().Header.Get("content-type") != "application/json" {
				return echo.NewHTTPError(http.StatusUnsupportedMediaType, "we only accept JSON data, sorry.")
			}
			// For valid credentials call next
			return next(ctx)
		}
	})

	// Routes
	apis.GET("", showAllItems(db))
	apis.GET("/:id", showOneItem(db))
	apis.POST("", createItem(db))
	apis.PUT("/:id", updateItem(db))
	apis.DELETE("/:id", deleteOneItem(db))
	apis.DELETE("", deleteManyItems(db))

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", *options.Port)))
}
