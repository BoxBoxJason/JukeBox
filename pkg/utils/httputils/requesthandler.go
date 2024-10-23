package httputils

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

// RetrieveAuthorizationToken retrieves the authorization token from the request header.
func RetrieveAuthorizationToken(r *http.Request, authorization_scheme string) (string, error) {
	auth_header := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth_header, authorization_scheme) {
		return "", NewUnauthorizedError("Invalid authorization scheme")
	}
	return strings.TrimPrefix(auth_header, authorization_scheme), nil
}

// RetrieveChiStringArgument retrieves a string argument from the URL parameters.
func RetrieveChiStringArgument(r *http.Request, argument_name string) (string, error) {
	argument := chi.URLParam(r, argument_name)
	if argument == "" {
		return "", NewBadRequestError("Missing argument: " + argument_name)
	}
	return argument, nil
}

// RetrieveChiIntArgument retrieves an integer argument from the URL parameters.
func RetrieveChiIntArgument(r *http.Request, argument_name string) (int, error) {
	argument := chi.URLParam(r, argument_name)
	if argument == "" {
		return -1, NewBadRequestError("Missing argument: " + argument_name)
	}
	return strconv.Atoi(argument)
}

// RetrieveStringParameter retrieves a string field from the form data.
// and returns an error if there are multiple values for the same key.
func RetrieveStringParameter(r *http.Request, field_name string, missing_ok bool) (string, error) {
	err := r.ParseForm()
	if err != nil {
		return "", NewBadRequestError("Failed to parse form")
	}

	values := r.Form[field_name]
	if len(values) == 0 && !missing_ok {
		return "", NewBadRequestError("Missing parameter: " + field_name)
	}
	if len(values) > 1 {
		return "", NewBadRequestError("Multiple values found for parameter: " + field_name)
	}
	return values[0], nil
}

// RetrieveIntParameter retrieves a single-value parameter from form data
// and returns an error if there are multiple values for the same key.
func RetrieveIntParameter(r *http.Request, parameter_name string, missing_ok bool) (string, error) {
	err := r.ParseForm()
	if err != nil {
		return "", NewBadRequestError("Failed to parse form data")
	}

	values := r.Form[parameter_name]
	if len(values) == 0 && !missing_ok {
		return "", NewBadRequestError("Missing parameter: " + parameter_name)
	}
	if len(values) > 1 {
		return "", NewBadRequestError("Multiple values found for parameter: " + parameter_name)
	}
	return values[0], nil
}

// RetrieveStringListValueParameter retrieves a list of values for a given parameter name from form data.
func RetrieveStringListValueParameter(r *http.Request, parameter_name string, missing_ok bool) ([]string, error) {
	if err := r.ParseForm(); err != nil {
		return nil, NewBadRequestError("Failed to parse form data")
	}
	values := r.Form[parameter_name]
	if len(values) == 0 && !missing_ok {
		return nil, NewBadRequestError("Missing parameter: " + parameter_name)
	}
	return values, nil
}

// RetrieveIntListValueParameter retrieves a list of integer values for a given parameter name from form data.
func RetrieveIntListValueParameter(r *http.Request, parameter_name string, missing_ok bool) ([]int, error) {
	if err := r.ParseForm(); err != nil {
		return nil, NewBadRequestError("Failed to parse form data")
	}
	values := r.Form[parameter_name]
	if len(values) == 0 && !missing_ok {
		return nil, NewBadRequestError("Missing parameter: " + parameter_name)
	}

	int_values := make([]int, 0, len(values))
	for _, value := range values {
		int_value, err := strconv.Atoi(value)
		if err != nil {
			return nil, NewBadRequestError("Invalid integer value: " + value)
		}
		int_values = append(int_values, int_value)
	}
	return int_values, nil
}

// RetrievePostFormStringParameter retrieves a single-value parameter from form data
func RetrievePostFormIntParameter(r *http.Request, parameter_name string, missing_ok bool) (int, error) {
	err := r.ParseForm()
	if err != nil {
		return -1, NewBadRequestError("Failed to parse form data")
	}

	values := r.PostForm[parameter_name]
	if len(values) == 0 && !missing_ok {
		return -1, NewBadRequestError("Missing parameter: " + parameter_name)
	}
	if len(values) > 1 {
		return -1, NewBadRequestError("Multiple values found for parameter: " + parameter_name)
	}

	return strconv.Atoi(values[0])
}

// RetrievePostFormStringParameter retrieves a single-value parameter from form data
func RetrievePostFormStringParameter(r *http.Request, parameter_name string, missing_ok bool) (string, error) {
	err := r.ParseForm()
	if err != nil {
		return "", NewBadRequestError("Failed to parse form data")
	}

	values := r.PostForm[parameter_name]
	if len(values) == 0 {
		if !missing_ok {
			return "", NewBadRequestError("Missing parameter: " + parameter_name)
		}
		return "", nil
	}
	if len(values) > 1 {
		return "", NewBadRequestError("Multiple values found for parameter: " + parameter_name)
	}

	return values[0], nil
}

// RetrievePostFormStringListValueParameter retrieves a list of values for a given parameter name from form data.
func RetrievePostFormStringListValueParameter(r *http.Request, parameter_name string, missing_ok bool) ([]string, error) {
	if err := r.ParseForm(); err != nil {
		return nil, NewBadRequestError("Failed to parse form data")
	}
	values := r.PostForm[parameter_name]
	if len(values) == 0 && !missing_ok {
		return nil, NewBadRequestError("Missing parameter: " + parameter_name)
	}
	return values, nil
}

// RetrievePostFormIntListValueParameter retrieves a list of integer values for a given parameter name from form data.
func RetrievePostFormIntListValueParameter(r *http.Request, parameter_name string, missing_ok bool) ([]int, error) {
	if err := r.ParseForm(); err != nil {
		return nil, NewBadRequestError("Failed to parse form data")
	}
	values := r.PostForm[parameter_name]
	if len(values) == 0 && !missing_ok {
		return nil, NewBadRequestError("Missing parameter: " + parameter_name)
	}

	int_values := make([]int, 0, len(values))
	for _, value := range values {
		int_value, err := strconv.Atoi(value)
		if err != nil {
			return nil, NewBadRequestError("Invalid integer value: " + value)
		}
		int_values = append(int_values, int_value)
	}
	return int_values, nil
}

func ReadCookie(r *http.Request, cookie_name string) (string, error) {
	cookie, err := r.Cookie(cookie_name)
	if err != nil {
		return "", NewUnauthorizedError("Missing cookie: " + cookie_name)
	}
	return cookie.Value, nil
}
