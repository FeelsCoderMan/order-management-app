package dto

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=5,max=30,alphanum"`
	Email    string `json:"email"    validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=12,max=50"`
}
