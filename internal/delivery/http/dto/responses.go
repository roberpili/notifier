package dto

import "github.com/rober0xf/notifier/internal/domain/entity"

// users
type LoginResponse struct {
	ID    int    `json:"id" example:"1"`
	Email string `json:"email" example:"rober@example.com"`
	Role  string `json:"role" example:"user"`
}

type UserCreatedResponse struct {
	Message string      `json:"message" example:"check your email to verify your account"`
	User    UserPayload `json:"user"`
}

type UserValidationErrorResponse struct {
	Error      string `json:"error"      example:"invalid email domain"`
	Suggestion string `json:"suggestion" example:"did you mean @gmail.com?"`
}

type VerifyEmailResponse struct {
	Message string `json:"message" example:"email verified successfully"`
}

// payments
type PaymentValidationErrorResponse struct {
	Error string         `json:"error"      example:"invalid category"`
	Meta  ValidationMeta `json:"meta"`
}

type ValidationMeta struct {
	AllowedTypes       []entity.TransactionType `json:"allowed_types"       example:"expense,income,subscription"`
	AllowedCategories  []entity.CategoryType    `json:"allowed_categories"  example:"electronics,entertainment,education,clothing,work,sports"`
	AllowedFrequencies []entity.FrequencyType   `json:"allowed_frequencies" example:"daily,weekly,monthly,yearly"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}
