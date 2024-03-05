package main

import (
	"go-api/internal/handlers"
	"go-api/internal/middleware"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	// Apply middleware to all routes
	r.Use(middleware.AuthMiddleware())

	// Protected routes
	r.GET("/", handlers.HomeHandler)
	r.GET("/user/:name", handlers.GetUserHandler)
	r.GET("/albums", handlers.GetAlbumsHandler)

	// Accessible routes
	r.GET("/auth/github/login", middleware.AuthMiddleware())
	r.GET("/auth/github/callback", middleware.AuthMiddleware())

	return r
}

func main() {
	r := setupRouter()
	log.Fatal(http.ListenAndServe(":8080", r))
}
