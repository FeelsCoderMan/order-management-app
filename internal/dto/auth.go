package dto

type SuccessLoginResponse struct {
	Success bool       `default:"true" json:"success"`
	AccessToken string `json:"accessToken"`
	ExpiresAt string   `json:"expiresAt"`
}

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=12,max=50"`
}
