package dto

type ErrorResponse struct {
	Success bool     `default:"false" json:"success"`
	Message []string `json:"message"`
}

type SuccessResponse struct {
	Success bool   `default:"true" json:"success"`
	Message string `json:"message"`
}

type SuccessLoginResponse struct {
	Success bool       `default:"true" json:"success"`
	AccessToken string `json:"accessToken"`

}
