package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)


type Validator struct{
	validator *validator.Validate
}

func (cv *Validator) Validate(i interface{}) error {
    if err := cv.validator.Struct(i); err != nil {
        return err
    }
    return nil
}

func RegisterValidator(router *echo.Echo){
	router.Validator = &Validator{validator: validator.New()}
}