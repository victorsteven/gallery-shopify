package utils

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
	"strings"
)

func ValidateInputs(dataSet interface{}) error {

	var validate *validator.Validate

	validate = validator.New()

	err := validate.Struct(dataSet)

	if err != nil {
		//Validation syntax is invalid
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return errors.New("validation syntax is invalid")
		}

		reflected := reflect.ValueOf(dataSet)

		for _, err := range err.(validator.ValidationErrors) {

			// Attempt to find field by name and get json tag name
			field, _ := reflected.Type().FieldByName(err.StructField())

			//If json tag doesn't exist, use lower case of name
			name := field.Tag.Get("json")
			if name == "" {
				name = strings.ToLower(err.StructField())
			}

			switch err.Tag() {
			case "required":
				return errors.New(fmt.Sprintf("The %s is required", name))
			case "email":
				return errors.New(fmt.Sprintf("The %s should be a valid email", name))
			default:
				return errors.New(fmt.Sprintf("The %s is invalid", name))
			}
		}
	}

	return nil
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
