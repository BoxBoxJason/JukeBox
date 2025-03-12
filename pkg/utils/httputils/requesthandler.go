package httputils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

// RetrieveAuthorizationToken retrieves the authorization token from the request header.
func RetrieveAuthorizationToken(r *http.Request, authorization_scheme string) (string, error) {
	auth_header := r.Header.Get("Authorization")
	if auth_header == "" {
		return "", NewUnauthorizedError("missing authorization header")
	} else if !strings.HasPrefix(auth_header, authorization_scheme) {
		return "", NewUnauthorizedError("invalid authorization scheme")
	}
	return strings.TrimSpace(strings.TrimPrefix(auth_header, authorization_scheme)), nil
}

// RetrieveChiStringArgument retrieves a string argument from the URL parameters.
func RetrieveChiStringArgument(r *http.Request, argument_name string) (string, error) {
	argument := chi.URLParam(r, argument_name)
	if argument == "" {
		return "", NewBadRequestError("missing argument: " + argument_name)
	}
	return strings.TrimSpace(argument), nil
}

// RetrieveChiIntArgument retrieves an integer argument from the URL parameters.
func RetrieveChiIntArgument(r *http.Request, argument_name string) (int, error) {
	argument := strings.TrimSpace(chi.URLParam(r, argument_name))
	if argument == "" {
		return 0, NewBadRequestError("missing argument: " + argument_name)
	}
	return strconv.Atoi(argument)
}

// RetrieveStringParameter retrieves a string field from the form data.
// and returns an error if there are multiple values for the same key.
func RetrieveStringParameter(r *http.Request, field_name string, missing_ok bool) (string, error) {
	err := r.ParseForm()
	if err != nil {
		return "", NewBadRequestError("failed to parse form")
	}

	values := r.Form[field_name]
	if len(values) == 0 {
		if !missing_ok {
			return "", NewBadRequestError("missing parameter: " + field_name)
		} else {
			return "", nil
		}
	}
	if len(values) > 1 {
		return "", NewBadRequestError("multiple values found for parameter: " + field_name)
	}
	return strings.TrimSpace(values[0]), nil
}

// RetrieveTimeStampParameter retrieves a single-value parameter from form data
// and returns an error if there are multiple values for the same key.
func RetrieveTimeStampParameter(r *http.Request, parameter_name string, missing_ok bool) (time.Time, error) {
	err := r.ParseForm()
	if err != nil {
		return time.Time{}, NewBadRequestError("failed to parse form data")
	}

	values := r.Form[parameter_name]
	if len(values) == 0 {
		if !missing_ok {
			return time.Time{}, NewBadRequestError("missing parameter: " + parameter_name)
		}
		return time.Time{}, nil
	}
	if len(values) > 1 {
		return time.Time{}, NewBadRequestError("multiple values found for parameter: " + parameter_name)
	}

	timestamp, err := time.Parse(time.RFC3339, strings.TrimSpace(values[0]))
	if err != nil {
		return time.Time{}, NewBadRequestError("invalid timestamp format: " + values[0])
	}
	return timestamp, nil
}

// RetrieveIntParameter retrieves a single-value parameter from form data
// and returns an error if there are multiple values for the same key.
func RetrieveIntParameter(r *http.Request, parameter_name string, missing_ok bool) (int, error) {
	err := r.ParseForm()
	if err != nil {
		return 0, NewBadRequestError("failed to parse form data")
	}

	values := r.Form[parameter_name]
	if len(values) == 0 {
		if !missing_ok {
			return 0, NewBadRequestError("missing parameter: " + parameter_name)
		}
		return 0, nil
	}
	if len(values) > 1 {
		return 0, NewBadRequestError("multiple values found for parameter: " + parameter_name)
	}
	return strconv.Atoi(strings.TrimSpace(values[0]))
}

// RetrieveBoolParameter retrieves a boolean value for a given parameter name from form data.
func RetrieveBoolParameter(r *http.Request, parameter_name string, missing_ok bool) ([]bool, error) {
	err := r.ParseForm()
	if err != nil {
		return []bool{}, NewBadRequestError("failed to parse form data")
	}

	values := r.Form[parameter_name]
	if len(values) == 0 {
		if !missing_ok {
			return []bool{}, NewBadRequestError("missing parameter: " + parameter_name)
		} else {
			return []bool{}, nil
		}
	}
	if len(values) > 1 {
		return []bool{}, NewBadRequestError("multiple values found for parameter: " + parameter_name)
	}
	result, err := strconv.ParseBool(strings.TrimSpace(values[0]))
	if err != nil {
		return []bool{}, NewBadRequestError("invalid boolean value: " + values[0])
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
	for i, value := range values {
		values[i] = strings.TrimSpace(value)
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
		int_value, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil {
			return nil, NewBadRequestError("Invalid integer value: " + value)
		}
		int_values = append(int_values, int_value)
	}
	return int_values, nil
}

// RetrievePostFormStringParameter retrieves a single-value parameter from form data
func RetrievePostFormIntParameter(r *http.Request, parameter_name string, missing_ok bool) (int, error) {
	contentType := r.Header.Get("Content-Type")

	if contentType == "application/x-www-form-urlencoded" {
		err := r.ParseForm()
		if err != nil {
			return 0, NewBadRequestError("Failed to parse form data")
		}
		values := r.PostForm[parameter_name]
		if len(values) == 0 {
			if !missing_ok {
				return 0, NewBadRequestError("Missing parameter: " + parameter_name)
			}
			return 0, nil
		}
		return strconv.Atoi(strings.TrimSpace(values[0]))
	} else if strings.Contains(contentType, "application/json") {
		// Buffer the body to allow re-reading
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			return 0, NewBadRequestError("Failed to read request body")
		}
		// Restore the original body for further use
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Decode the JSON payload
		var body map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &body); err != nil {
			return 0, NewBadRequestError("Invalid JSON payload")
		}

		// Extract the parameter
		if value, exists := body[parameter_name]; exists {
			if floatValue, ok := value.(float64); ok { // JSON numbers are float64
				return int(floatValue), nil
			}
			return 0, NewBadRequestError("Invalid type for parameter: " + parameter_name)
		}
		if !missing_ok {
			return 0, NewBadRequestError("Missing parameter: " + parameter_name)
		}
	}

	return 0, NewBadRequestError("Unsupported Content-Type")
}

// RetrievePostFormStringParameter retrieves a single-value parameter from form data
func RetrievePostFormStringParameter(r *http.Request, parameter_name string, missing_ok bool) (string, error) {
	contentType := r.Header.Get("Content-Type")

	if contentType == "application/x-www-form-urlencoded" {
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
		return strings.TrimSpace(values[0]), nil
	} else if strings.Contains(contentType, "application/json") {
		// Buffer the body to allow re-reading
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			return "", NewBadRequestError("Failed to read request body")
		}
		// Restore the original body for further use
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Decode the JSON payload
		var body map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &body); err != nil {
			return "", NewBadRequestError("Invalid JSON payload")
		}

		// Extract the parameter
		if value, exists := body[parameter_name]; exists {
			if str, ok := value.(string); ok {
				return strings.TrimSpace(str), nil
			}
			return "", NewBadRequestError("Invalid type for parameter: " + parameter_name)
		}

		if !missing_ok {
			return "", NewBadRequestError("Missing parameter: " + parameter_name)
		}
	}

	return "", NewBadRequestError("Unsupported Content-Type")
}

// RetrievePostFormBoolParamater retrieves a boolean value for a given parameter name from form data.
func RetrievePostFormBoolParameter(r *http.Request, parameter_name string, missing_ok bool) ([]bool, error) {
	contentType := r.Header.Get("Content-Type")

	if contentType == "application/x-www-form-urlencoded" {
		err := r.ParseForm()
		if err != nil {
			return nil, NewBadRequestError("Failed to parse form data")
		}
		values := r.PostForm[parameter_name]
		if len(values) == 0 {
			if !missing_ok {
				return nil, NewBadRequestError("Missing parameter: " + parameter_name)
			}
			return nil, nil
		}
		result, err := strconv.ParseBool(strings.TrimSpace(values[0]))
		if err != nil {
			return nil, NewBadRequestError("Invalid boolean value: " + values[0])
		}
		return []bool{result}, nil
	} else if strings.Contains(contentType, "application/json") {
		// Buffer the body to allow re-reading
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, NewBadRequestError("Failed to read request body")
		}
		// Restore the original body for further use
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Decode the JSON payload
		var body map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &body); err != nil {
			return nil, NewBadRequestError("Invalid JSON payload")
		}

		// Extract the parameter
		if value, exists := body[parameter_name]; exists {
			if boolValue, ok := value.(bool); ok {
				return []bool{boolValue}, nil
			}
			return nil, NewBadRequestError("Invalid type for parameter: " + parameter_name)
		}
		if !missing_ok {
			return nil, NewBadRequestError("Missing parameter: " + parameter_name)
		}
	}

	return nil, NewBadRequestError("Unsupported Content-Type")
}

// RetrievePostFormStringListValueParameter retrieves a list of values for a given parameter name from form data.
func RetrievePostFormStringListValueParameter(r *http.Request, parameter_name string, missing_ok bool) ([]string, error) {
	contentType := r.Header.Get("Content-Type")

	if contentType == "application/x-www-form-urlencoded" {
		if err := r.ParseForm(); err != nil {
			return nil, NewBadRequestError("Failed to parse form data")
		}
		values := r.PostForm[parameter_name]
		if len(values) == 0 && !missing_ok {
			return nil, NewBadRequestError("Missing parameter: " + parameter_name)
		}
		for i, value := range values {
			values[i] = strings.TrimSpace(value)
		}
		return values, nil
	} else if strings.Contains(contentType, "application/json") {
		// Buffer the body to allow re-reading
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, NewBadRequestError("Failed to read request body")
		}
		// Restore the original body for further use
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Decode the JSON payload
		var body map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &body); err != nil {
			return nil, NewBadRequestError("Invalid JSON payload")
		}

		// Extract the parameter
		if value, exists := body[parameter_name]; exists {
			if list, ok := value.([]interface{}); ok {
				strList := []string{}
				for _, item := range list {
					if str, ok := item.(string); ok {
						strList = append(strList, strings.TrimSpace(str))
					} else {
						return nil, NewBadRequestError("Invalid type in list for parameter: " + parameter_name)
					}
				}
				return strList, nil
			}
			return nil, NewBadRequestError("Invalid type for parameter: " + parameter_name)
		}
		if !missing_ok {
			return nil, NewBadRequestError("Missing parameter: " + parameter_name)
		}
	}

	return nil, NewBadRequestError("Unsupported Content-Type")
}

// RetrievePostFormIntListValueParameter retrieves a list of integer values for a given parameter name from form data.
func RetrievePostFormIntListValueParameter(r *http.Request, parameter_name string, missing_ok bool) ([]int, error) {
	contentType := r.Header.Get("Content-Type")

	if contentType == "application/x-www-form-urlencoded" {
		if err := r.ParseForm(); err != nil {
			return nil, NewBadRequestError("Failed to parse form data")
		}
		values := r.PostForm[parameter_name]
		if len(values) == 0 && !missing_ok {
			return nil, NewBadRequestError("Missing parameter: " + parameter_name)
		}

		intValues := make([]int, len(values))
		for i, value := range values {
			intValue, err := strconv.Atoi(strings.TrimSpace(value))
			if err != nil {
				return nil, NewBadRequestError("Invalid integer value: " + value)
			}
			intValues[i] = intValue
		}
		return intValues, nil
	} else if strings.Contains(contentType, "application/json") {
		// Buffer the body to allow re-reading
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, NewBadRequestError("Failed to read request body")
		}
		// Restore the original body for further use
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		// Decode the JSON payload
		var body map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &body); err != nil {
			return nil, NewBadRequestError("Invalid JSON payload")
		}
		// Extract the parameter
		if value, exists := body[parameter_name]; exists {
			if list, ok := value.([]interface{}); ok {
				intList := []int{}
				for _, item := range list {
					if floatValue, ok := item.(float64); ok {
						intList = append(intList, int(floatValue))
					} else {
						return nil, NewBadRequestError("Invalid type in list for parameter: " + parameter_name)
					}
				}
				return intList, nil
			}
			return nil, NewBadRequestError("Invalid type for parameter: " + parameter_name)
		}
		if !missing_ok {
			return nil, NewBadRequestError("Missing parameter: " + parameter_name)
		}
	}

	return nil, NewBadRequestError("Unsupported Content-Type")
}

func ReadCookie(r *http.Request, cookie_name string) (string, error) {
	cookie, err := r.Cookie(cookie_name)
	if err != nil {
		return "", errors.New("missing cookie: " + cookie_name)
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

func RetryRequest(request_url string, body []byte, retry_delay int, max_retry int) error {
	if max_retry <= 0 {
		max_retry = 1
	}

	for i := 0; i < max_retry; i++ {
		resp, err := http.Post(request_url, "application/json", bytes.NewBuffer(body))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			return nil
		}

		time.Sleep(time.Duration(retry_delay))
	}

	return errors.New("failed to send prompt after maximum retries")
}
