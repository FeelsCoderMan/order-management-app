package main

import (
	"net/http"

	"github.com/FeelsCoderMan/order-management-app/internal/config"
	"github.com/FeelsCoderMan/order-management-app/internal/database"
	"github.com/FeelsCoderMan/order-management-app/internal/dto"
	"github.com/FeelsCoderMan/order-management-app/internal/handler"
	"github.com/FeelsCoderMan/order-management-app/internal/middleware"
	"github.com/FeelsCoderMan/order-management-app/internal/model"
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

	authConfig, err := config.LoadAuthConfig()

	if err != nil {
		mainLogger.Fatal(err)
	}

	registerLogger := utils.GetLogger("[Register] ")
	authMiddlewareLogger := utils.GetLogger("[AuthMiddleware] ")

	userRepository := repository.NewUserRepository(db)
	authService := service.NewAuthService(
		authConfig,
		userRepository,
		registerLogger,
	)
	authHandler := handler.NewAuthHandler(
		authService,
		registerLogger,
	)
	authMiddleware := middleware.NewAuthMiddleware(
		authConfig,
		authMiddlewareLogger,
		userRepository,
	)

	router := mux.NewRouter()
	publicRouter := router.PathPrefix("/").Subrouter()
	// TODO: Add private subrouter to use auth middleware
	privateRouter := router.PathPrefix("/").Subrouter()
	// TODO: Add restricted HTTP Methods for each handler
	publicRouter.HandleFunc("/register", authHandler.Register)
	publicRouter.HandleFunc("/login", authHandler.Login)
	// router.HandleFunc("/logout", authHandler.Logout)
	privateRouter.Use(authMiddleware.RequireAuthenticate)
	// TODO: Remove test handler as it was used for testing purposes
	privateRouter.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		ctxValue := r.Context().Value(middleware.UserKey)

		if user, ok := ctxValue.(*model.User); ok {
			utils.WriteJSON(w, http.StatusOK, &dto.SuccessResponse{
				Message: user.ID,
			})
			return
		}

		utils.WriteJSON(w, http.StatusBadRequest, &dto.ErrorResponse{
			Message: []string{"User does not exist"},
		})
	})

	mainLogger.Println("Server running on :8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		mainLogger.Fatal(err)
	}
}
