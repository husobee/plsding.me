package main

import (
	"database/sql"
	"log"

	"github.com/husobee/plsding.me/handlers"
	"github.com/husobee/plsding.me/models"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// create a new echo instance
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Route / to handler function
	e.GET("/health-check", handlers.HealthCheck)

	// Authenticated Routes
	var signingKey = []byte("superdupersecret!")
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(models.SigningContextKey, signingKey)
			return next(c)
		}
	})

	// add database to context
	var db, err = sql.Open("sqlite3", "./plsding.me.db")
	if err != nil {
		log.Fatalf("error opening database: %v\n", err)
	}
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(models.DBContextKey, db)
			return next(c)
		}
	})

	// Authentication routes
	e.POST("/login", handlers.Login)
	e.POST("/logout", handlers.Logout)

	g := e.Group("/reminder")
	g.Use(middleware.JWT(signingKey))
	g.POST("", handlers.CreateReminder)

	// start the server, and log if it fails
	e.Logger.Fatal(e.Start(":8080"))
}
