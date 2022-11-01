package main

import (
	"github.com/Oxyrus/shopping/internal/controllers"
	"github.com/Oxyrus/shopping/internal/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sqlx.Open("postgres", "postgres://andres:@localhost/shopping?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r := gin.Default()

	api := r.Group("/api")

	userController := controllers.UserController{DB: db}
	api.POST("/login", userController.Login)
	api.POST("/register", userController.Register)

	protected := api.Group("/profile")
	protected.Use(middlewares.AuthMiddleware())

	healthController := controllers.HealthController{}
	protected.GET("/health", healthController.GetHealthStatus)

	if err = r.Run(); err != nil {
		panic(err)
	}
}
