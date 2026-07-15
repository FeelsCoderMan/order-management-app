package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/FeelsCoderMan/order-management-app/internal/dto"
	"github.com/FeelsCoderMan/order-management-app/internal/service"
	"github.com/FeelsCoderMan/order-management-app/internal/utils"
	"github.com/FeelsCoderMan/order-management-app/internal/validation"
)


type AuthHandler struct {
	authService service.AuthService
	logger *log.Logger
}

func NewAuthHandler(authService service.AuthService, logger *log.Logger ) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger: logger,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var userReq dto.CreateUserRequest;

	if err := utils.FromJSON(r, &userReq); err != nil {
		h.logger.Println("Failed parsing register request: ", err)

		utils.WriteJSON(w, http.StatusBadRequest, dto.ErrorResponse{
			Message: []string{"Invalid request body"},
		})
		return
	}

	if errs := utils.ValidateStruct(h.logger, userReq); len(errs) > 0 {
		utils.WriteJSON(w, http.StatusBadRequest, dto.ErrorResponse{
			Message: errs,
		})
		return
	}

	if err := validation.ValidateUsername(userReq.Username); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, dto.ErrorResponse{
			Message: []string{err.Error()},
		})
		return
	}

	// TODO: Add context.Context support for request timeout and cancellation handling.
	if err := h.authService.Register(userReq); err != nil {
		switch {
		case errors.Is(err, service.ErrEmailAlreadyUsed):
			utils.WriteJSON(w, http.StatusConflict, dto.ErrorResponse{
				Message: []string{"Email is already used"},
			})
		default:
			h.logger.Println("Register user failed: ", err)
			utils.WriteJSON(w, http.StatusInternalServerError, dto.ErrorResponse{
				Message: []string{"Something went wrong"},
			})
		}

		return
	}

	utils.WriteJSON(w, http.StatusCreated, dto.SuccessResponse{
		Message: "User registered successfully.",
	});
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginReq dto.LoginRequest;

	if err := utils.FromJSON(r, &loginReq); err != nil {
		h.logger.Println("Failed to parse login request: ", err)

		utils.WriteJSON(w, http.StatusBadRequest, dto.ErrorResponse{
			Message: []string{"Invalid request body"},
		})
		return
	}

	if errs := utils.ValidateStruct(h.logger, loginReq); len(errs) > 0 {
		utils.WriteJSON(w, http.StatusBadRequest, dto.ErrorResponse{
			Message: errs,
		})
		return
	}

	accessToken, err := h.authService.Login(loginReq);

	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			utils.WriteJSON(w, http.StatusUnauthorized, dto.ErrorResponse{
				Message: []string{"Invalid email or password"},
			})
		default:
			utils.WriteJSON(w, http.StatusInternalServerError, dto.ErrorResponse{
				Message: []string{"Something went wrong"},
			})
		}

		return
	}

	utils.WriteJSON(w, http.StatusOK, dto.SuccessLoginResponse{
		AccessToken: accessToken,
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
}
