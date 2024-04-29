package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/seemsod1/ancy/internal/config"
	"github.com/seemsod1/ancy/internal/handlers"
	"net/http"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Logger)

	mux.Get("/user-role", handlers.Repo.GetAllUserRoles)
	mux.Post("/user-role/create", handlers.Repo.CreateUserRole)
	mux.Put("/user-role/update", handlers.Repo.UpdateUserRole)
	mux.Delete("/user-role/delete/{id}", handlers.Repo.DeleteUserRole)

	mux.Get("/user", handlers.Repo.GetAllUsers)
	mux.Post("/user/create", handlers.Repo.CreateUser)
	//mux.Put("/user/update", handlers.Repo.UpdateUser)
	//mux.Delete("/user/delete/{id}", handlers.Repo.DeleteUser)

	join := chi.NewRouter()
	join.Use(SessionLoad)

	//join.Get("/singUp", controllers.Repo.SingUpPage)
	//join.Post("/singUp", controllers.Repo.UserSingUp)
	//
	join.Post("/logout", handlers.Repo.Logout)
	join.Post("/login", handlers.Repo.Login)
	join.Post("/sing-up", handlers.Repo.SingUp)
	mux.Mount("/join", join)

	return mux
}
