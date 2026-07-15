package main

import (
	"net/http"

	"github.com/FeelsCoderMan/order-management-app/internal/database"
	"github.com/FeelsCoderMan/order-management-app/internal/handler"
	"github.com/FeelsCoderMan/order-management-app/internal/repository"
	"github.com/FeelsCoderMan/order-management-app/internal/service"
	"github.com/FeelsCoderMan/order-management-app/internal/utils"
	"github.com/gorilla/mux"
)

func main() {
	mainLogger := utils.GetLogger("[main] ")
	db, err := database.InitPostgres()

	if err != nil {
		mainLogger.Fatal(err)
	}

	registerLogger := utils.GetLogger("[Register] ")
	userRepository := repository.NewUserRepository(
		db,
	)

	authService := service.NewAuthService(
		userRepository,
		registerLogger,
	)

	authHandler := handler.NewAuthHandler(
		authService,
		registerLogger,
	)

	router := mux.NewRouter()
	router.HandleFunc("/register", authHandler.Register)
	router.HandleFunc("/login", authHandler.Login)
	// router.HandleFunc("/logout", authHandler.Logout)

	mainLogger.Println("Server running on :8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		mainLogger.Fatal(err)
	}
}
