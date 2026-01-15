package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	"tourbackend/internal/auth"
	"tourbackend/internal/courses"
	"tourbackend/internal/courses/materials"
	"tourbackend/internal/courses/quizzes"
	db "tourbackend/internal/database"
	"tourbackend/internal/feeds"
	"tourbackend/internal/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/joho/godotenv"
)

// this variable changes some stuff based on whether the app is deployed or not
// for example the cookie will be set to secure - requiring https
var IS_DEPLOYED bool

// if this variable is true, then the db gets deleted and seeded each reload of the server
var RESET_DB bool

var STATIC_PATH string = "../../static"

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

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

	IS_DEPLOYED = strings.ToLower(os.Getenv("IS_DEPLOYED")) == "true"
	RESET_DB = strings.ToLower(os.Getenv("RESET_DB")) == "true"
	SEED := strings.ToLower(os.Getenv("SEED")) == "true"

	db, queries := db.Initialize(RESET_DB)
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
	e.POST("/logout", authHandler.Logout)

	//* Course Feeds
	feedsService := feeds.NewService(queries, "./static")
	feedsHandler := feeds.NewHandler(STATIC_PATH, feedsService, queries, IS_DEPLOYED)

	e.GET("/courses/:courseId/feed", feedsHandler.GetCourseFeed)
	e.POST("/courses/:courseId/feed", feedsHandler.CreateFeedPost)
	e.PUT("/courses/:courseId/feed/:postId", feedsHandler.UpdateFeedPost)
	e.DELETE("/courses/:courseId/feed/:postId", feedsHandler.DeleteFeedPost)

	e.GET("/courses/:courseId/feed/stream", feedsHandler.StreamFeed)

	//* Courses and it's deps (materials and quizzes - TODO)
	matsService := materials.NewService(queries, STATIC_PATH, feedsService)
	quizzesService := quizzes.NewService(queries, STATIC_PATH, feedsService)

	courseService := courses.NewService(queries, matsService, quizzesService, feedsService)

	coursesHandler := courses.NewCourseHandler(queries, IS_DEPLOYED, courseService)

	e.GET("/courses", coursesHandler.ListAllCourses)
	e.POST("/courses", coursesHandler.CreateCourse)

	e.GET("/courses/:courseId", coursesHandler.GetCourse)
	e.PUT("/courses/:courseId", coursesHandler.UpdateCourse)
	e.DELETE("/courses/:courseId", coursesHandler.DeleteCourse)

	//* Course materials
	materialsHandler := materials.NewHandler(STATIC_PATH, matsService, queries, IS_DEPLOYED)

	materials := e.Group("/courses/:courseId/materials")
	materials.GET("", materialsHandler.ListMaterials)
	materials.POST("", materialsHandler.CreateMaterial)

	materials.PUT("/:materialId", materialsHandler.UpdateMaterial)
	materials.DELETE("/:materialId", materialsHandler.DeleteMaterial)

	//* Course Quizes
	quizzesHandler := quizzes.NewHandler(STATIC_PATH, quizzesService, queries, IS_DEPLOYED)

	quizzes := e.Group("/courses/:courseId/quizzes")
	quizzes.GET("", quizzesHandler.ListQuizzes)
	quizzes.POST("", quizzesHandler.CreateQuiz)

	quizzes.GET("/:quizId", quizzesHandler.GetQuiz)
	quizzes.PUT("/:quizId", quizzesHandler.UpdateQuiz)
	quizzes.DELETE("/:quizId", quizzesHandler.DeleteQuiz)

	quizzes.POST("/:quizId/submit", quizzesHandler.SubmitQuizAnswers)

	//* Static
	e.Static("/static", STATIC_PATH)

	// Seed the db with 3 courses
	if SEED {
		err := seed(queries, courseService, feedsService, matsService, quizzesService)
		if err != nil {
			panic(err)
		}
	}

	// Create the admin account as described in the 1. phase
	err = createAdmin(queries)
	if err != nil {
		panic(err)
	}

	fmt.Println("ready!")

	e.Logger.Fatal(e.Start(":" + PORT_STRING))
}
