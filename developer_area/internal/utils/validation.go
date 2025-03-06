package utils

import (
	"context"

	"encore.dev/beta/errs"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/mold/v4/modifiers"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type ValidationErrors map[string][]string

func (e *ValidationErrors) ErrDetails() {}

func ValidateTransform(ctx context.Context, params any) error {
	eb := errs.B().Code(errs.InvalidArgument)

	validate := validator.New()
	conform := modifiers.New()

	err := conform.Struct(ctx, params)
	if err != nil {
		return err
	}

	err = validate.Struct(params)
	if err != nil {
		details := GetValidationErrorDetails(validate, err)

		eb.Msg("Validation error").Details(&details)

		return eb.Err()
	}

	return err
}

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
