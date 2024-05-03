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
	mux.Use(SessionLoad)

	// Роути для управління ролями користувачів
	mux.Route("/user-role", func(r chi.Router) {
		r.Get("/", handlers.Repo.GetAllUserRoles)              // Тільки для адміна
		r.Post("/create", handlers.Repo.CreateUserRole)        // Тільки для адміна
		r.Put("/update", handlers.Repo.UpdateUserRole)         // Тільки для адміна
		r.Delete("/delete/{id}", handlers.Repo.DeleteUserRole) // Тільки для адміна
	})

	// Роутер для залогінених користувачів
	authRouter := chi.NewRouter()
	authRouter.Use(SessionLoad)
	authRouter.Use(AuthUser)

	authRouter.Post("/logout", handlers.Repo.Logout)
	authRouter.Post("/exhibit/create", handlers.Repo.CreateExhibit) // Зареєстровані користувачі
	//authRouter.Put("/exhibit/update", handlers.Repo.UpdateExhibit) // Зареєстровані користувачі
	authRouter.Delete("/exhibit/delete/{id}", handlers.Repo.DeleteExhibit) // Зареєстровані користувачі
	authRouter.Get("/exhibit/my", handlers.Repo.GetMyExhibits)             // Зареєстровані користувачі
	authRouter.Get("/me", handlers.Repo.GetMe)                             // Зареєстровані користувачі

	// Роутер для адміністратора
	adminRouter := chi.NewRouter()
	adminRouter.Use(SessionLoad)
	adminRouter.Use(AuthUser)
	adminRouter.Use(AdminOnly)

	adminRouter.Post("/exhibit/approve/{id}", handlers.Repo.ApproveExhibit) // Тільки для адміна

	mux.Mount("/user", authRouter)   // Встановлюємо роутер для залогінених користувачів
	mux.Mount("/admin", adminRouter) // Встановлюємо роутер для адміністратора

	// Роутер для гостя
	mux.Post("/login", handlers.Repo.Login)            // Гість
	mux.Post("/sign-up", handlers.Repo.SignUp)         // Гість
	mux.Get("/exhibit", handlers.Repo.GetAllExhibits)  // Гість
	mux.Get("/exhibit/{id}", handlers.Repo.GetExhibit) // Гість

	return mux
}
