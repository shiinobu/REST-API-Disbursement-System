package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("gt0", func(fl validator.FieldLevel) bool {
			value := fl.Field().String()
			amount, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return false
			}

			return amount > 0
		})
	}
}

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
	case "gt0":
		return "harus lebih besar dari 0"
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
