package validator

import (
	"errors"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/validator/v10"

	ut "github.com/go-playground/universal-translator"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type Validator struct {
	validate *validator.Validate
	trans    ut.Translator
}

func NewValidator() *Validator {
	validate := validator.New()

	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, trans)

	val := Validator{
		validate: validate,
		trans:    trans,
	}

	return &val
}

func (v *Validator) Validate(i interface{}) error {
	if err := v.validate.Struct(i); err != nil {
		var errs = []string{}

		validatorErrs := err.(validator.ValidationErrors)
		for _, e := range validatorErrs {
			errs = append(errs, e.Translate(v.trans))
		}

		return errors.New(strings.Join(errs, ", "))
	}
	return nil
}
