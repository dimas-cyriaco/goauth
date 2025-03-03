package utils

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type ValidationErrors map[string][]string

func (e *ValidationErrors) ErrDetails() {}

func GetValidationErrorDetails(validate *validator.Validate, err error) ValidationErrors {
	validationErrors := err.(validator.ValidationErrors)

	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")

	en_translations.RegisterDefaultTranslations(validate, trans)

	translationsMap := make(ValidationErrors)

	for _, e := range validationErrors {
		translation := e.Translate(trans)
		field := ToSnakeCase(e.Field())

		translationsMap[field] = append(translationsMap[field], translation)
	}

	return translationsMap
}
