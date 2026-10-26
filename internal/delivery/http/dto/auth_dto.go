package dto

type RegisterPayload struct {
	Username string `json:"username" binding:"required,min=4,max=32" example:"rober"`
	Email    string `json:"email" binding:"required,email,max=254" example:"rober@example.com"`
	Password string `json:"password" binding:"required,min=8,max=72" example:"secret123#!"`
}

type LoginPayload struct {
	Email    string `json:"email" binding:"required,email,max=254" example:"rober@example.com"`
	Password string `json:"password" binding:"required,min=8,max=72" example:"secret123#!"`
}

type GoogleLoginRequest struct {
	IDToken string `json:"id_token" binding:"required" example:"eyGaksdjaldasdj..."`
}
