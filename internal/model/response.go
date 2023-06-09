package model

import "UsersService/internal/cErrors"

type ResponseStandard struct {
	Success *ResponseSuccessModel       `json:"response"`
	Error   *cErrors.ResponseErrorModel `json:"error"`
}

type ResponseSuccessModel struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
