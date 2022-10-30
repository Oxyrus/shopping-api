package main

import (
	"github.com/Oxyrus/shopping/internal/controllers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"net/http"
)

const secret = "<secret>"

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte(secret), nil)
}

func main() {
	db, err := sqlx.Open("postgres", "postgres://andres:@localhost/shopping?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r := chi.NewRouter()

	r.Use(middleware.AllowContentEncoding("application/json"))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		healthController := controllers.HealthController{}
		r.Get("/health", healthController.GetHealthStatus)

		userController := controllers.UserController{
			TokenAuth: tokenAuth,
			DB:        db,
		}
		r.Post("/login", userController.Login)
		r.Post("/register", userController.Register)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))

			r.Get("/profile", userController.Profile)
		})
	})

	err = http.ListenAndServe(":4000", r)
	if err != nil {
		panic(err)
	}
}
