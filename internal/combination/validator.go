package combination

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

func Validator() (*validator.Validate, error) {
	var validate = validator.New()

	// register a custom validation for combination data
	err := validate.RegisterValidation("combination_data_json", func(fl validator.FieldLevel) bool {
		value, isRightType := fl.Field().Interface().(map[DataType]*Data)
		if !isRightType {
			return false
		}

		for _, data := range value {
			if data.Type == DataTypeJSON {
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
