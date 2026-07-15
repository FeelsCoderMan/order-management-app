package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
)

func FromJSON(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func WriteJSON[T any](w http.ResponseWriter, status int, v *T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	logger := GetLogger("[utils-json] ")

	if v == nil {
		logger.Println("WriteJSON called with nil response body")
		return
	}

	if err := defaults.Set(v); err != nil {
		logger.Printf("Could not apply defaults: %v\n", err)
		return
	}

	if err := json.NewEncoder(w).Encode(v); err != nil {
		logger.Printf("Could not write json response: %v\n", err)
	}
}

func ValidateStruct(logger *log.Logger, v any) []string {
	validate := validator.New()

	if err := validate.Struct(v); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := make([]string, 0, len(validationErrors))

			for _, e := range validationErrors {
				switch e.Tag() {
					case "required":
						errors = append(errors, fmt.Sprintf("%s is required", e.Field()))
					case "min":
						errors = append(errors, fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param()))
					case "max":
						errors = append(errors, fmt.Sprintf("%s must be at most %s characters", e.Field(), e.Param()))
					case "alphanum":
						errors = append(errors, fmt.Sprintf("%s must contain only alphanumerical characters", e.Field()))
					case "email":
						errors = append(errors, fmt.Sprintf("%s must be a valid email address", e.Field()))
					default:
						logger.Printf("%s case is not handled. Default case is triggered\n", e.Tag())
						errors = append(errors, fmt.Sprintf("%s is invalid", e.Field()))
				}
			}

			return errors;
		}
	}

	return nil
}

func GetLogger(prefix string) *log.Logger {
	return log.New(os.Stdout, prefix, log.LstdFlags)
}
