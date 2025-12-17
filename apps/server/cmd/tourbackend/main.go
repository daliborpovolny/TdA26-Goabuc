package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"tourbackend/internal/auth"
	"tourbackend/internal/courses"
	"tourbackend/internal/courses/materials"
	db "tourbackend/internal/database"
	"tourbackend/internal/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// this variable changes some stuff based on whether the app is deployed or not
// for example the cookie will be set to secure - requiring https
var IS_DEPLOYED bool

// if this variable is true, then the db gets deleted and seeded each reload of the server
var RESET_DB bool

var STATIC_PATH string = "../../static"

func main() {

	// in future setting env vars should not be done here
	os.Setenv("IS_DEPLOYED", "false")
	os.Setenv("RESET_DB", "true")

	// try to read the port number from env, if fails default to 3000
	PORT_STRING := os.Getenv("PORT")
	_, err := strconv.Atoi(PORT_STRING)
	if err != nil || PORT_STRING == "" {
		PORT_STRING = "3000"
	}

	// if env var sets the value of /static path -> respect it (this makes it work both locally and in docker)
	ENV_STATIC_PATH := os.Getenv("STATIC_PATH")
	if ENV_STATIC_PATH != "" {
		STATIC_PATH = ENV_STATIC_PATH
	}

	IS_DEPLOYED = os.Getenv("IS_DEPLOYED") == "true"
	RESET_DB = os.Getenv("RESET_DB") == "true"
	RESET_DB = false

	db, queries := db.Initialize(RESET_DB, RESET_DB)
	defer db.Close()
	fmt.Println("initialized db")

	e := echo.New()
	// e.Debug = true // enabling this make echo log more stuff into the console

	slog.SetDefault(middlewares.Logger)
	e.Use(middlewares.LoggerMiddleware)

	e.Use(middleware.Recover())
	e.Use(auth.AuthMiddleware(queries))

	e.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"organization": "Student Cyber Games"})
	})

	//* Auth
	authHandler := auth.NewAuthHandler(queries, IS_DEPLOYED)

	e.POST("/register", authHandler.Register)
	e.POST("/login", authHandler.Login)
	e.GET("/me", authHandler.Profile)

	//* Courses and it's deps (materials and quizzes - TODO)
	matsService := materials.NewService(queries, STATIC_PATH)
	courseService := courses.NewService(queries, matsService)

	coursesHandler := courses.NewCourseHandler(queries, IS_DEPLOYED, courseService)

	e.GET("/courses", coursesHandler.ListAllCourses)
	e.POST("/courses", coursesHandler.CreateCourse)

	e.GET("/courses/:courseId", coursesHandler.GetCourse)
	e.PUT("/courses/:courseId", coursesHandler.UpdateCourse)
	e.DELETE("/courses/:courseId", coursesHandler.DeleteCourse)

	//* Course materials
	materialsHandler := materials.NewHandler(STATIC_PATH, queries, IS_DEPLOYED, matsService)

	materials := e.Group("/courses/:courseId/materials")
	materials.GET("", materialsHandler.ListMaterials)
	materials.POST("", materialsHandler.CreateMaterial)

	materials.PUT("/:materialId", materialsHandler.UpdateMaterial)
	materials.DELETE("/:materialId", materialsHandler.DeleteMaterial)

	//* Static
	e.Static("/static", STATIC_PATH)

	fmt.Println("ready!")

	e.Logger.Fatal(e.Start(":" + PORT_STRING))
}
