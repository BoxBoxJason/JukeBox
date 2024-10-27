package httputils

import (
	"mime/multipart"
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
	return strings.TrimSpace(strings.TrimPrefix(auth_header, authorization_scheme)), nil
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
		return 0, NewBadRequestError("Missing argument: " + argument_name)
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
	if len(values) == 0 {
		if !missing_ok {
			return "", NewBadRequestError("Missing parameter: " + field_name)
		} else {
			return "", nil
		}
	}
	if len(values) > 1 {
		return "", NewBadRequestError("Multiple values found for parameter: " + field_name)
	}
	return values[0], nil
}

// RetrieveIntParameter retrieves a single-value parameter from form data
// and returns an error if there are multiple values for the same key.
func RetrieveIntParameter(r *http.Request, parameter_name string, missing_ok bool) (int, error) {
	err := r.ParseForm()
	if err != nil {
		return 0, NewBadRequestError("Failed to parse form data")
	}

	values := r.Form[parameter_name]
	if len(values) == 0 {
		if !missing_ok {
			return 0, NewBadRequestError("Missing parameter: " + parameter_name)
		}
		return 0, nil
	}
	if len(values) > 1 {
		return 0, NewBadRequestError("Multiple values found for parameter: " + parameter_name)
	}
	return strconv.Atoi(values[0])
}

// RetrieveBoolParameter retrieves a boolean value for a given parameter name from form data.
func RetrieveBoolParameter(r *http.Request, parameter_name string, missing_ok bool) ([]bool, error) {
	err := r.ParseForm()
	if err != nil {
		return []bool{}, NewBadRequestError("Failed to parse form data")
	}

	values := r.Form[parameter_name]
	if len(values) == 0 {
		if !missing_ok {
			return []bool{}, NewBadRequestError("Missing parameter: " + parameter_name)
		} else {
			return []bool{}, nil
		}
	}
	if len(values) > 1 {
		return []bool{}, NewBadRequestError("Multiple values found for parameter: " + parameter_name)
	}
	result, err := strconv.ParseBool(values[0])
	if err != nil {
		return []bool{}, NewBadRequestError("Invalid boolean value: " + values[0])
	}
	return []bool{result}, nil
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
		return 0, NewBadRequestError("Failed to parse form data")
	}

	values := r.PostForm[parameter_name]
	if len(values) == 0 {
		if !missing_ok {
			return 0, NewBadRequestError("Missing parameter: " + parameter_name)
		} else {
			return 0, nil
		}
	}
	if len(values) > 1 {
		return 0, NewBadRequestError("Multiple values found for parameter: " + parameter_name)
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
		} else {
			return "", nil
		}
	}
	if len(values) > 1 {
		return "", NewBadRequestError("Multiple values found for parameter: " + parameter_name)
	}

	return values[0], nil
}

// RetrievePostFormBoolParamater retrieves a boolean value for a given parameter name from form data.
func RetrievePostFormBoolParameter(r *http.Request, parameter_name string, missing_ok bool) ([]bool, error) {
	err := r.ParseForm()
	if err != nil {
		return []bool{}, NewBadRequestError("Failed to parse form data")
	}

	values := r.PostForm[parameter_name]
	if len(values) == 0 {
		if !missing_ok {
			return []bool{}, NewBadRequestError("Missing parameter: " + parameter_name)
		}
		return []bool{}, nil
	}
	if len(values) > 1 {
		return []bool{}, NewBadRequestError("Multiple values found for parameter: " + parameter_name)
	}

	result, err := strconv.ParseBool(values[0])
	if err != nil {
		return []bool{}, NewBadRequestError("Invalid boolean value: " + values[0])
	}
	return []bool{result}, nil
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

// RetrieveImageFile retrieves an uploaded image file from a multipart form.
// It returns the file and its header or an error if the file is missing or not accessible.
func RetrieveImageFile(r *http.Request, fieldName string, missingOk bool) (multipart.File, *multipart.FileHeader, error) {
	// Ensure the request size is under 3MB
	err := r.ParseMultipartForm(3 << 20)
	if err != nil {
		return nil, nil, NewBadRequestError("Request size exceeds the limit")
	}
	// Retrieve the file from the specified field name
	file, file_header, err := r.FormFile(fieldName)
	if err != nil {
		if err == http.ErrMissingFile && missingOk {
			// If missing_ok is true, return nil without raising an error
			return nil, nil, nil
		}
		return nil, nil, NewBadRequestError("Missing or inaccessible file: " + fieldName)
	}

	// Check if the uploaded file is an image
	err = ValidateImageFile(file, file_header)
	if err != nil {
		return nil, nil, err
	}

	return file, file_header, nil
}

// ValidateImageFile checks if the uploaded file is an image by verifying the MIME type.
// It returns an error if the file is not an image.
func ValidateImageFile(file multipart.File, file_header *multipart.FileHeader) error {
	// Check if fileHeader's content type starts with "image/"
	if !strings.HasPrefix(file_header.Header.Get("Content-Type"), "image/") {
		return NewBadRequestError("Uploaded file is not an image")
	}

	// Reset the file pointer to the beginning of the file
	file.Seek(0, 0)

	// Read the first 512 bytes to verify image file type
	buf := make([]byte, 512)
	_, err := file.Read(buf)
	if err != nil {
		return NewBadRequestError("Failed to read the file")
	}

	// Reset the file pointer to the beginning of the file
	file.Seek(0, 0)

	// Check if the content type of the file is indeed an image type by inspecting the buffer
	contentType := http.DetectContentType(buf)
	if !strings.HasPrefix(contentType, "image/") {
		return NewBadRequestError("Uploaded file does not appear to be a valid image")
	}

	return nil
}
