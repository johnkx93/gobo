package validation

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// Validator wraps go-playground validator with English translations
type Validator struct {
	validate *validator.Validate
	trans    ut.Translator
}

// New creates and configures a new Validator with English translations
func New() *Validator {
	validate := validator.New()

	// Register custom tag name function to use JSON field names in error messages
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "" || name == "-" {
			return fld.Name
		}
		return name
	})

	// Setup English locale
	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)

	// Register English translations
	trans, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, trans)

	v := &Validator{
		validate: validate,
		trans:    trans,
	}

	// Register custom validators
	v.registerCustomValidators()

	// Register custom translations/messages
	v.registerCustomMessages()

	return v
}

// Struct validates a struct and returns error
func (v *Validator) Struct(s interface{}) error {
	return v.validate.Struct(s)
}

// GetValidator returns the underlying validator instance
func (v *Validator) GetValidator() *validator.Validate {
	return v.validate
}

// TranslateErrors translates the first validation error to English
func (v *Validator) TranslateErrors(err error) string {
	if err == nil {
		return ""
	}

	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err.Error()
	}

	// If no validation errors, return empty string
	if len(validationErrs) == 0 {
		return ""
	}

	// Translate only the first error
	firstErr := validationErrs[0]
	translatedMsg := firstErr.Translate(v.trans)

	return translatedMsg
}

// registerCustomValidators registers custom validation functions
func (v *Validator) registerCustomValidators() {
	// Example: Custom validator to check if string doesn't contain bad words
	v.validate.RegisterValidation("notbadword", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		badWords := []string{"badword", "spam", "offensive"}

		valueLower := strings.ToLower(value)
		for _, word := range badWords {
			if strings.Contains(valueLower, word) {
				return false
			}
		}
		return true
	})

	// Example: Custom validator for strong password
	v.validate.RegisterValidation("strongpassword", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()

		// Check length
		if len(password) < 8 {
			return false
		}

		// Check for at least one uppercase, one lowercase, and one digit
		hasUpper := false
		hasLower := false
		hasDigit := false

		for _, char := range password {
			switch {
			case char >= 'A' && char <= 'Z':
				hasUpper = true
			case char >= 'a' && char <= 'z':
				hasLower = true
			case char >= '0' && char <= '9':
				hasDigit = true
			}
		}

		return hasUpper && hasLower && hasDigit
	})
}

// registerCustomMessages registers custom error messages for validation tags
func (v *Validator) registerCustomMessages() {
	// Custom message for "required" tag in English
	v.validate.RegisterTranslation("required", v.trans,
		func(ut ut.Translator) error {
			return ut.Add("required", "{0} is required and cannot be empty", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("required", fe.Field())
			return t
		},
	)

	// Custom message for "email" tag in English
	v.validate.RegisterTranslation("email", v.trans,
		func(ut ut.Translator) error {
			return ut.Add("email", "{0} must be a valid email address", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("email", fe.Field())
			return t
		},
	)

	// Custom message for "min" tag in English
	v.validate.RegisterTranslation("min", v.trans,
		func(ut ut.Translator) error {
			return ut.Add("min", "{0} must be at least {1} characters long", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("min", fe.Field(), fe.Param())
			return t
		},
	)

	// Register custom validator messages
	v.validate.RegisterTranslation("notbadword", v.trans,
		func(ut ut.Translator) error {
			return ut.Add("notbadword", "{0} contains inappropriate content", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("notbadword", fe.Field())
			return t
		},
	)

	v.validate.RegisterTranslation("strongpassword", v.trans,
		func(ut ut.Translator) error {
			return ut.Add("strongpassword", "{0} must be at least 8 characters and contain uppercase, lowercase, and digit", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("strongpassword", fe.Field())
			return t
		},
	)
}
