package main

import (
	"database/sql"

	"github.com/husobee/plsding.me/handlers"
	"github.com/husobee/plsding.me/middlewares"
	"github.com/husobee/plsding.me/models"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// create a new echo instance
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)

	e.Pre(middlewares.RequestIDMiddleware)

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

	e.Static("/static/*", "static")

	// V1 Routes
	v1 := e.Group("/v1")
	// V1 Authentication routes
	v1.POST("/login", handlers.Login)
	v1.POST("/logout", handlers.Logout)
	// V1 Reminder Routes
	v1Reminders := v1.Group("/reminder", middleware.JWT(signingKey))
	v1Reminders.POST("", handlers.CreateReminder)
	v1Reminders.GET("/:id", handlers.GetReminder)

	// Latest Authentication routes
	e.POST("/login", handlers.Login)
	e.POST("/logout", handlers.Logout)

	// Latest Reminder Routes
	g := e.Group("/reminder")
	g.Use(middleware.JWT(signingKey))
	g.POST("", handlers.CreateReminder)
	g.GET(":id", handlers.GetReminder)

	// start the server, and log if it fails
	e.Logger.Fatal(e.Start(":8080"))
}
