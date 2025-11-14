package main

import (
	"fmt"
	"os"

	"net/http"
	db "tourbackend/database"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// this variable changes some stuff based on whether the app is deployed or not
// for example the cookie will be set to secure - requiring https
var IS_DEPLOYED bool

// if this variable is true, then the db gets deleted and seeded each reload of the server
var RESET_DB bool

func main() {

	// in future setting env vars should not be done here
	os.Setenv("IS_DEPLOYED", "false")
	os.Setenv("RESET_DB", "true")

	IS_DEPLOYED = os.Getenv("IS_DEPLOYED") == "true"
	RESET_DB = os.Getenv("RESET_DB") == "true"

	db, queries := db.Initialize()
	defer db.Close()
	fmt.Println("initialized db")

	e := echo.New()
	e.Debug = true // enabling this make echo log more stuff into the console

	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(AuthMiddleware(queries))

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"organization": "Student Cyber Games"})
	})

	authHandler := NewAuthHandler(queries)

	e.POST("/register", authHandler.register)
	e.POST("/login", authHandler.login)
	e.GET("/me", authHandler.profile)

	e.Static("/static", "static")

	e.Logger.Fatal(e.Start(":1323"))
}
