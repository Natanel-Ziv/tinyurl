package utils

import (
	"net/mail"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var validateEmail validator.Func = func(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	_, err := mail.ParseAddress(email)
	return err == nil
}

func InitValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validateemail", validateEmail)
	}
}
