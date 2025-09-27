package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationMiddleware struct {
	validator *validator.Validate
}

type ValidationError struct {
	Type    string   `json:"type"` // "schema" or "validation"
	Message string   `json:"message"`
	Errors  []string `json:"errors,omitempty"`
}

func NewValidationMiddleware() *ValidationMiddleware {
	return &ValidationMiddleware{
		validator: validator.New(),
	}
}

func ValidateBody[T any](vm *ValidationMiddleware, dto T) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				vm.writeValidationError(w, "schema", "Unable to read request body", nil)
				return
			}

			if len(body) == 0 {
				vm.writeValidationError(w, "schema", "Request body is required", nil)
				return
			}

			var rawData map[string]any
			if err := json.Unmarshal(body, &rawData); err != nil {
				vm.writeValidationError(w, "schema", "Invalid JSON format", []string{err.Error()})
				return
			}

			var validatedData T

			if err := json.Unmarshal(body, &validatedData); err != nil {
				vm.writeValidationError(w, "schema", "JSON structure doesn't match expected format", []string{err.Error()})
				return
			}

			if unknownFields := vm.checkUnknownFields(rawData, validatedData); len(unknownFields) > 0 {
				vm.writeValidationError(w, "schema", "Unknown fields detected", unknownFields)
				return
			}

			if err := vm.validator.Struct(validatedData); err != nil {
				validationErrors := vm.formatValidationErrors(err)
				vm.writeValidationError(w, "validation", "Field validation failed", validationErrors)
				return
			}

			ctx := context.WithValue(r.Context(), "validatedData", validatedData)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (vm *ValidationMiddleware) checkUnknownFields(rawData map[string]any, dtoValue any) []string {
	var unknownFields []string

	validFields := vm.getStructJSONFields(dtoValue)

	for field := range rawData {
		if _, exists := validFields[field]; !exists {
			unknownFields = append(unknownFields, fmt.Sprintf("'%s' is not a valid field", field))
		}
	}

	return unknownFields
}

func (vm *ValidationMiddleware) getStructJSONFields(dto any) map[string]bool {
	fields := make(map[string]bool)

	dtoType := reflect.TypeOf(dto)
	if dtoType.Kind() == reflect.Ptr {
		dtoType = dtoType.Elem()
	}

	for i := 0; i < dtoType.NumField(); i++ {
		field := dtoType.Field(i)
		jsonTag := field.Tag.Get("json")

		if jsonTag == "" {
			fields[strings.ToLower(field.Name)] = true
		} else if jsonTag != "-" {
			jsonName := strings.Split(jsonTag, ",")[0]
			if jsonName != "" {
				fields[jsonName] = true
			}
		}
	}

	return fields
}

func (vm *ValidationMiddleware) writeValidationError(w http.ResponseWriter, errorType, message string, errors []string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	validationErr := ValidationError{
		Type:    errorType,
		Message: message,
		Errors:  errors,
	}

	json.NewEncoder(w).Encode(validationErr)
}

func (vm *ValidationMiddleware) formatValidationErrors(err error) []string {
	var errors []string

	for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			errors = append(errors, fmt.Sprintf("'%s' is required", err.Field()))
		case "email":
			errors = append(errors, fmt.Sprintf("'%s' must be a valid email", err.Field()))
		case "min":
			errors = append(errors, fmt.Sprintf("'%s' must be at least %s characters", err.Field(), err.Param()))
		case "max":
			errors = append(errors, fmt.Sprintf("'%s' must be at most %s characters", err.Field(), err.Param()))
		case "url":
			errors = append(errors, fmt.Sprintf("'%s' must be a valid URL", err.Field()))
		case "oneof":
			errors = append(errors, fmt.Sprintf("'%s' must be one of: %s", err.Field(), err.Param()))
		case "gte":
			errors = append(errors, fmt.Sprintf("'%s' must be greater than or equal to %s", err.Field(), err.Param()))
		case "lte":
			errors = append(errors, fmt.Sprintf("'%s' must be less than or equal to %s", err.Field(), err.Param()))
		case "gt":
			errors = append(errors, fmt.Sprintf("'%s' must be greater than %s", err.Field(), err.Param()))
		case "lt":
			errors = append(errors, fmt.Sprintf("'%s' must be less than %s", err.Field(), err.Param()))
		default:
			errors = append(errors, fmt.Sprintf("'%s' failed validation: %s", err.Field(), err.Tag()))
		}
	}

	return errors
}

func GetValidatedData[T any](r *http.Request) T {
	data := r.Context().Value("validatedData")
	return data.(T)
}
