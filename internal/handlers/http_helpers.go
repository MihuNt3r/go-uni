package handlers

import (
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

type errorResponse struct {
	Error string `json:"error"`
}

type messageResponse struct {
	Message string `json:"message"`
}

func parsePositiveID(raw string) (int64, error) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, err
	}
	if id <= 0 {
		return 0, errors.New("id must be positive")
	}

	return id, nil
}

func parsePathID(r *http.Request, name string) (int64, error) {
	raw := r.PathValue(name)
	if raw == "" {
		return 0, errors.New("missing path parameter")
	}

	return parsePositiveID(raw)
}

func validatePayload(payload any) error {
	err := Validate.Struct(payload)
	if err == nil {
		return nil
	}

	var validationErrs validator.ValidationErrors
	if !errors.As(err, &validationErrs) || len(validationErrs) == 0 {
		return errors.New("invalid request body")
	}

	firstErr := validationErrs[0]
	fieldName := jsonFieldName(payload, firstErr.StructField())

	switch firstErr.Tag() {
	case "required":
		return errors.New(fieldName + " is required")
	case "email":
		return errors.New(fieldName + " must be valid")
	case "gt":
		return errors.New(fieldName + " must be positive")
	default:
		return errors.New(fieldName + " is invalid")
	}
}

func jsonFieldName(payload any, structFieldName string) string {
	t := reflect.TypeOf(payload)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	field, ok := t.FieldByName(structFieldName)
	if !ok {
		return strings.ToLower(structFieldName)
	}

	jsonTag := field.Tag.Get("json")
	if jsonTag == "" || jsonTag == "-" {
		return strings.ToLower(structFieldName)
	}

	parts := strings.Split(jsonTag, ",")
	if parts[0] == "" {
		return strings.ToLower(structFieldName)
	}

	return parts[0]
}
