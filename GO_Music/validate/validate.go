package validate

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
)

var (
	Validate     *validator.Validate
	ValidateOnce sync.Once
)

// Инициализация валидатора и регистрация кастомных валидаторов
func InitValidator() {
	Validate = validator.New()

	// Регистрация кастомных валидаторов
	Validate.RegisterValidation("birthday_past", BirthdayPast)
	Validate.RegisterValidation("role_valid", RoleValid)
	Validate.RegisterValidation("grade_range", GradeRange)
}

// Получение валидатора с одинарной инициализацией
func GetValidator() *validator.Validate {
	ValidateOnce.Do(InitValidator)
	return Validate
}

// Гарантированное использование библиотеки для предотвращения её удаления при tidy
var _ = GetValidator()

type ValidationErrors map[string]string

func (ve ValidationErrors) Error() string {
	var msg string
	for field, err := range ve {
		msg += fmt.Sprintf("%s: %s; ", field, err)
	}
	return msg
}

func ValidateStruct(s interface{}) error {
	v := GetValidator()
	err := v.Struct(s)
	if err == nil {
		return nil
	}

	if _, ok := err.(*validator.InvalidValidationError); ok {
		return fmt.Errorf("invalid validation error: %w", err)
	}

	validationErrors := ValidationErrors{}
	for _, err := range err.(validator.ValidationErrors) {
		field := err.Field()
		tag := err.Tag()
		var message string

		switch tag {
		case "required":
			message = "is required"
		case "min":
			message = fmt.Sprintf("minimum value is %s", err.Param())
		case "max":
			message = fmt.Sprintf("maximum value is %s", err.Param())
		case "email":
			message = "must be a valid email"
		case "e164":
			message = "must be a valid E.164 phone number"
		case "birthday_past":
			message = "date must not be in the future"
		case "role_valid":
			message = "role is not valid"
		case "grade_range":
			message = "grade must be between 0 and 100"
		default:
			message = fmt.Sprintf("failed on '%s' validation", tag)
		}
		validationErrors[field] = message
	}
	return validationErrors
}

// --- Кастомные валидаторы ---

// BirthdayPast проверяет, что дата рождения не в будущем
func BirthdayPast(fl validator.FieldLevel) bool {
	birthday, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}
	now := time.Now()
	return birthday.Before(now) || birthday.Equal(now)
}

// RoleValid проверяет, что роль входит в список допустимых
func RoleValid(fl validator.FieldLevel) bool {
	role, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	allowedRoles := []string{"student", "employee", "admin", "teacher"}
	role = strings.ToLower(role)
	for _, r := range allowedRoles {
		if role == r {
			return true
		}
	}
	return false
}

// GradeRange проверяет, что оценка в диапазоне 0..10
func GradeRange(fl validator.FieldLevel) bool {
	grade, ok := fl.Field().Interface().(int)
	if !ok {
		return false
	}
	return grade >= 0 && grade <= 10
}
