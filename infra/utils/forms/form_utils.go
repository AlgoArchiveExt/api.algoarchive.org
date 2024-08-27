package formutils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Gets the field name from a form given a form property that has a json tag.
func getJSONFieldNameFromFormProperty(form interface{}, property string) (value string, ok bool) {
	t := reflect.TypeOf(form)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return "Form was not an interface or struct.", false
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if field.Name == property {
			jsonFieldName, found := field.Tag.Lookup("json")
			if !found {
				return "", false
			}

			return jsonFieldName, true
		}

		if field.Type.Kind() == reflect.Struct {
			tryFindPropertyInDeeperStruct, found := getJSONFieldNameFromFormProperty(reflect.New(field.Type).Interface(), property)
			if found {
				thisFieldJSONName, _ := getJSONFieldNameFromFormProperty(form, field.Name)
				return fmt.Sprintf("%s.%s", thisFieldJSONName, tryFindPropertyInDeeperStruct), true
			}
		}
	}

	return "", false
}

func getErrorsFromForm(err error, form interface{}) []string {
	errors := []string{}

	for _, e := range err.(validator.ValidationErrors) {
		jsonFieldName, ok := getJSONFieldNameFromFormProperty(form, e.Field())

		if ok {
			errors = append(errors, jsonFieldName)
		}
	}

	return errors
}

// GenerateJSONBindingErrorMessage generates an error message for HTTP requests when binding its body to a form.
// It expects a form (which should be a pointer) and its error from the binding attempt.
// Usually, this error is a [validator.ValidationError], meaning that the request body is missing some required form fields.
func GenerateJSONBindingErrorMessage(form interface{}, err error) (message string) {
	switch err.(type) {
	case validator.ValidationErrors:
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return "Some fields are invalid or missing"
		}

		var errors []string = getErrorsFromForm(err, form)

		if len(errors) == 0 {
			formType := reflect.TypeOf(form)

			// The form argument should be a pointer already, but this check will be here just in case.
			if formType.Kind() == reflect.Ptr {
				formType = formType.Elem()
			}

			return
		}

		var errorString string = "The following fields are missing or invalid: " + strings.Join(errors, ", ")
		return errorString

	case *json.SyntaxError:
		return "You have a syntax error"

	default:
		return "Invalid request"
	}
}
