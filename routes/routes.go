package routes

import (
	"net/http"

	"auth_go/controllers"
	"auth_go/middleware"

	"github.com/gorilla/mux"
)

// SetupRoutes configures and returns a router with all application routes
func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Create instances of controllers
	userController := controllers.NewUserController()

	// Public routes
	router.HandleFunc("/api/register", middleware.LoggerMiddleware(userController.Register)).Methods("POST")
	router.HandleFunc("/api/login", middleware.LoggerMiddleware(userController.Login)).Methods("POST")

	// Protected routes
	router.HandleFunc("/api/profile", middleware.LoggerMiddleware(middleware.AuthMiddleware(userController.GetProfile))).Methods("GET")
	router.HandleFunc("/api/profile", middleware.LoggerMiddleware(middleware.AuthMiddleware(userController.UpdateProfile))).Methods("PUT")
	router.HandleFunc("/api/change-password", middleware.LoggerMiddleware(middleware.AuthMiddleware(userController.ChangePassword))).Methods("POST")

	// Apply CORS middleware to all routes
	handler := middleware.CORSMiddleware(router.ServeHTTP)

	return &mux.Router{
		SkipClean:      router.SkipClean,
		StrictSlash:    router.StrictSlash,
		UseEncodedPath: router.UseEncodedPath,
		Handler:        http.HandlerFunc(handler),
	}
}
