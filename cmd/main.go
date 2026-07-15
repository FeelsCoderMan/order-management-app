package main

import (
	"log"
	"net/http"
	"os"

	"github.com/FeelsCoderMan/order-management-app/internal/database"
	"github.com/FeelsCoderMan/order-management-app/internal/handler"
	"github.com/FeelsCoderMan/order-management-app/internal/repository"
	"github.com/FeelsCoderMan/order-management-app/internal/service"
	"github.com/gorilla/mux"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	db, err := database.InitPostgres()

	if err != nil {
		logger.Fatal(err)
	}

	userRepository := repository.NewUserRepository(
		db,
	)

	authService := service.NewAuthService(
		userRepository,
		logger,
	)

	authHandler := handler.NewAuthHandler(
		authService,
		logger,
	)

	router := mux.NewRouter()
	router.HandleFunc("/register", authHandler.Register)
	router.HandleFunc("/login", authHandler.Login)
	// router.HandleFunc("/logout", authHandler.Logout)

	logger.Println("Server running on :8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		logger.Fatal(err)
	}
}
