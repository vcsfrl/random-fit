package combination

import (
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

	return validate, err
}
