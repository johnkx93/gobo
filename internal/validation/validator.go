package validation

import (
	"net/http"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

// Validator wraps go-playground validator with i18n support
type Validator struct {
	validate *validator.Validate
	uni      *ut.UniversalTranslator
}

// New creates and configures a new Validator with i18n support for en, ms, zh
func New() *Validator {
	validate := validator.New()

	// Register custom tag name function to use JSON field names in error messages
	// validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
	// 	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	// 	if name == "" || name == "-" {
	// 		return fld.Name
	// 	}
	// 	return name
	// })

	// Setup locales
	enLocale := en.New()
	zhLocale := zh.New()

	// Create universal translator with English as fallback
	uni := ut.New(enLocale, enLocale, zhLocale)

	// Register default translations for supported languages
	enTrans, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, enTrans)

	zhTrans, _ := uni.GetTranslator("zh")
	zh_translations.RegisterDefaultTranslations(validate, zhTrans)

	v := &Validator{
		validate: validate,
		uni:      uni,
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

// TranslateErrors translates the first validation error based on Accept-Language header
func (v *Validator) TranslateErrors(r *http.Request, err error) string {
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

	// Get translator based on Accept-Language header
	locale := v.extractLocale(r)
	trans, found := v.uni.GetTranslator(locale)
	if !found {
		trans, _ = v.uni.GetTranslator("en") // fallback to English
	}

	// Translate only the first error
	firstErr := validationErrs[0]
	translatedMsg := firstErr.Translate(trans)

	return translatedMsg
}

// extractLocale extracts locale from Accept-Language header
// Returns "en", "ms", or "zh" based on header, defaults to "en"
func (v *Validator) extractLocale(r *http.Request) string {
	return "en"
	// acceptLang := r.Header.Get("Accept-Language")
	// if acceptLang == "" {
	// 	return "en"
	// }

	// // Parse Accept-Language header (simple implementation)
	// // Format: "en-US,en;q=0.9,ms;q=0.8,zh-CN;q=0.7"
	// parts := strings.Split(acceptLang, ",")
	// for _, part := range parts {
	// 	lang := strings.TrimSpace(strings.Split(part, ";")[0])
	// 	lang = strings.ToLower(lang)

	// 	// Match primary language code
	// 	if strings.HasPrefix(lang, "en") {
	// 		return "en"
	// 	}
	// 	if strings.HasPrefix(lang, "zh") {
	// 		return "zh"
	// 	}
	// }

	// return "en"
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
	enTrans, _ := v.uni.GetTranslator("en")
	v.validate.RegisterTranslation("required", enTrans,
		func(ut ut.Translator) error {
			return ut.Add("required", "{0} is required and cannot be empty", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("required", fe.Field())
			return t
		},
	)

	// Custom message for "email" tag in English
	v.validate.RegisterTranslation("email", enTrans,
		func(ut ut.Translator) error {
			return ut.Add("email", "{0} must be a valid email address", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("email", fe.Field())
			return t
		},
	)

	// Custom message for "min" tag in English
	v.validate.RegisterTranslation("min", enTrans,
		func(ut ut.Translator) error {
			return ut.Add("min", "{0} must be at least {1} characters long", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("min", fe.Field(), fe.Param())
			return t
		},
	)

	// Register custom validator messages
	v.validate.RegisterTranslation("notbadword", enTrans,
		func(ut ut.Translator) error {
			return ut.Add("notbadword", "{0} contains inappropriate content", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("notbadword", fe.Field())
			return t
		},
	)

	v.validate.RegisterTranslation("strongpassword", enTrans,
		func(ut ut.Translator) error {
			return ut.Add("strongpassword", "{0} must be at least 8 characters and contain uppercase, lowercase, and digit", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("strongpassword", fe.Field())
			return t
		},
	)

	// Chinese translations for custom validators
	zhTrans, _ := v.uni.GetTranslator("zh")
	v.validate.RegisterTranslation("notbadword", zhTrans,
		func(ut ut.Translator) error {
			return ut.Add("notbadword", "{0}包含不当内容", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("notbadword", fe.Field())
			return t
		},
	)

	v.validate.RegisterTranslation("strongpassword", zhTrans,
		func(ut ut.Translator) error {
			return ut.Add("strongpassword", "{0}必须至少8个字符，并包含大写、小写字母和数字", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("strongpassword", fe.Field())
			return t
		},
	)
}
