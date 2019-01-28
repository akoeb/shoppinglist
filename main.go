package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/acme/autocert"
)

// Options hold global application options that can be set via CLI
type Options struct {
	DatabaseFile         *string
	HTTPBaseAuthUser     *string
	HTTPBaseAuthPassword *string
	Domain               *string
	Port                 *int
	BindIP               *string
	LogLevel             log.Lvl
}

func parseFlags() *Options {

	options := &Options{}

	// options:
	options.Port = flag.Int("port", 8080, "The Port that the application uses for listening")
	options.Domain = flag.String("domain", "localhost", "The Domain for CORS and TLS")
	options.DatabaseFile = flag.String("db", "storage.db", "The file to store the sqlite3 database")
	options.HTTPBaseAuthUser = flag.String("user", "", "The user name for HTTP Base authentication")
	options.HTTPBaseAuthPassword = flag.String("password", "", "The password for HTTP Base authentication")
	options.BindIP = flag.String("bind", "", "The IP address to bind to, defaults to all local")
	debugFlag := flag.Bool("debug", false, "Activate debug logging")

	// parse command line into options
	flag.Parse()

	// some options need special treatment
	if *debugFlag {
		options.LogLevel = log.DEBUG
	} else {
		options.LogLevel = log.INFO
	}

	// CLI Validation goes here
	if *options.HTTPBaseAuthPassword == "" && *options.HTTPBaseAuthUser != "" {
		log.Fatal("Can not use HTTP Base Authentication with only user, needs also password")
	}
	if *options.HTTPBaseAuthPassword != "" && *options.HTTPBaseAuthUser == "" {
		log.Fatal("Can not use HTTP Base Authentication with only password, needs also user")
	}
	log.Infof("options: %v/%v", *options.HTTPBaseAuthUser, *options.HTTPBaseAuthPassword)
	return options
}

func main() {

	options := parseFlags()

	// Echo instance
	e := echo.New()
	e.Logger.SetLevel(options.LogLevel)

	// Database
	db := initDB(*options.DatabaseFile)

	// channel to send back and forth update notifications
	notifier := NewNotifier()

	// Middleware
	e.AutoTLSManager.Cache = autocert.DirCache(".cache")
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS header
	var corsConfig middleware.CORSConfig
	if *options.Domain == "localhost" {
		corsConfig = middleware.CORSConfig{
			AllowOrigins: []string{fmt.Sprintf("http://127.0.0.1:%d", *options.Port)},
			AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
		}
	} else {
		corsConfig = middleware.CORSConfig{
			AllowOrigins: []string{fmt.Sprintf("https://%s:%d", *options.Domain, *options.Port)},
			AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
		}
	}
	e.Use(middleware.CORSWithConfig(corsConfig))

	// basic auth
	if *options.HTTPBaseAuthUser != "" {
		e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
			if username == *options.HTTPBaseAuthUser && password == *options.HTTPBaseAuthPassword {
				return true, nil
			}
			return false, nil
		}))
	}

	// routes for static files
	e.Static("/", "/app/public/")

	// apis have their own middlewares: group them
	apis := e.Group("/api")

	// only allow application/json content type:
	apis.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if ctx.Request().Header.Get(echo.HeaderContentType) != echo.MIMEApplicationJSON {
				return echo.NewHTTPError(http.StatusUnsupportedMediaType, "we only accept JSON data, sorry.")
			}
			// For valid credentials call next
			return next(ctx)
		}
	})

	// TODO: notifier only for all categories or inside a category

	// Routes
	// categories
	apis.GET("/categories", showAllCategories(db))         // list all categories
	apis.GET("/categories/:catid", showOneCategory(db))    // show one category
	apis.POST("/categories", createCategory(db, notifier)) // create a category
	apis.PUT("/categories/:catid", updateCategory(db, notifier))
	apis.DELETE("/categories/:catid", deleteOneCategory(db, notifier))

	// locations
	apis.GET("/locations", showAllLocations(db))          // list all locations, can be filtered by category id
	apis.POST("/locations", createLocation(db, notifier)) // create a location
	apis.PUT("/locations/:locid", updateLocation(db, notifier))
	apis.DELETE("/locations/:locid", deleteOneLocation(db, notifier))
	// TODO: apis.POST("/locations/:locid/reorder", reorderItemsByLocation(db, notifier))

	// items
	apis.GET("/categories/:catid/items", showAllItemsInCategory(db)) // list all items in one category, can be filtered by location id
	apis.POST("/categories/:catid/items", createItem(db, notifier))  // create an item
	apis.GET("/categories/:catid/items/:itemid", showOneItemInCategory(db))
	apis.PUT("/categories/:catid/items/:itemid", updateItem(db, notifier))
	apis.DELETE("/categories/:catid/items/:itemid", deleteOneItem(db, notifier))

	// events
	events := e.Group("/events")
	events.GET("", eventsStream(notifier))

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", *options.BindIP, *options.Port)))

}
