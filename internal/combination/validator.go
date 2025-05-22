package combination

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

func Validator() (*validator.Validate, error) {
	var validate = validator.New()

	// register a custom validation for combination data
	err := validate.RegisterValidation("combination_data_json", func(fl validator.FieldLevel) bool {
		value := fl.Field().Interface().(map[DataType]*Data)

		for _, data := range value {
			if data.Type == DataTypeJson {
				return true
			}
		}

		return false
	})

	if err != nil {
		return nil, fmt.Errorf("register validator: %w", err)
	}

	return validate, nil
}
