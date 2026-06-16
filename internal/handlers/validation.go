package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidationError(c *gin.Context, err error, request any) {
	Error(c, http.StatusBadRequest, "Request tidak valid", validationMessages(err, request))
}

func validationMessages(err error, request any) gin.H {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return gin.H{"body": "format request tidak valid"}
	}

	messages := gin.H{}
	for _, fieldError := range validationErrors {
		field := jsonFieldName(request, fieldError.StructField())
		messages[field] = validationMessage(fieldError)
	}

	return messages
}

func validationMessage(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return "field wajib diisi"
	case "email":
		return "format email tidak valid"
	case "min":
		return fmt.Sprintf("minimal %s karakter", fieldError.Param())
	case "max":
		return fmt.Sprintf("maksimal %s karakter", fieldError.Param())
	case "gt":
		return fmt.Sprintf("harus lebih besar dari %s", fieldError.Param())
	case "numeric":
		return "hanya boleh berisi angka"
	default:
		return "nilai tidak valid"
	}
}

func jsonFieldName(request any, structField string) string {
	requestType := reflect.TypeOf(request)
	if requestType.Kind() == reflect.Pointer {
		requestType = requestType.Elem()
	}

	field, found := requestType.FieldByName(structField)
	if !found {
		return strings.ToLower(structField)
	}

	jsonName := strings.Split(field.Tag.Get("json"), ",")[0]
	if jsonName == "" || jsonName == "-" {
		return strings.ToLower(structField)
	}

	return jsonName
}
