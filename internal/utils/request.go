package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
)

func ParseBodyToJSON(w http.ResponseWriter, r *http.Request, model interface{}) (string, int, error) {
	// Limit request body size
	const maxBodySize = 10 * 1024 * 1024 // 10MB
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	// Create JSON decoder
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevent fields not in struct

	// Decode JSON body into struct
	if err := decoder.Decode(model); err != nil {
		var errMessage string
		var statusCode int

		switch {
		// Incomplete JSON data
		case err == io.ErrUnexpectedEOF:
			errMessage = "Request body is incomplete"
			statusCode = http.StatusBadRequest

		// Syntax error in the Request body
		case err.(*json.SyntaxError) != nil:
			errMessage = "Malformed JSON syntax"
			statusCode = http.StatusBadRequest

		// Incorrect types in the Request body
		case err.(*json.UnmarshalTypeError) != nil:
			errMessage = "Incorrect data type in JSON"
			statusCode = http.StatusBadRequest

		// Target struct passed directly
		case reflect.TypeOf(model).Kind() != reflect.Ptr:
			errMessage = "Target struct must be a pointer"
			statusCode = http.StatusInternalServerError

		// Default error
		default:
			errMessage = "Invalid request body"
			statusCode = http.StatusBadRequest
		}

		return errMessage, statusCode, err
	}

	return "", http.StatusOK, nil
}
