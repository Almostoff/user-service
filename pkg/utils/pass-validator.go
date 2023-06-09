package utils

import (
	"UsersService/internal/cErrors"
	"fmt"
	"strings"
	"unicode"
)

const symbols string = "/[!@#$%^&*(),.?\":{}|<>]/"

func ValidPassword(s string) (cErr cErrors.ResponseErrorModel) {
	var lengthEnough, number, lower, upper, special, isValid bool
	length := 0
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
			length++
		case unicode.IsUpper(c):
			upper = true
			length++
		case unicode.IsLower(c):
			lower = true
			length++
		case strings.Contains(symbols, fmt.Sprintf("%c", c)):
			special = true
			length++
		case unicode.IsLetter(c):
			length++
		default:
			length++
		}
	}
	lengthEnough = length >= 10
	isValid = lengthEnough && number && upper && special && lower
	if !isValid {
		fmt.Println(lengthEnough, number, upper, special, lower)
		cErr.InternalCode = cErrors.StatusBadRequest
		cErr.StandartCode = cErrors.StatusBadRequest
		cErr.Message = "password must contain: "
		if !number {
			cErr.Message += " number"
		}
		if !lengthEnough {
			cErr.Message += " ten or more symbols "
		}
		if !upper {
			cErr.Message += " letters in uppercase "
		}
		if !lower {
			cErr.Message += " letters in lowercase "
		}
		if !special {
			cErr.Message += " special symbols "
		}
	}
	return cErr
}
