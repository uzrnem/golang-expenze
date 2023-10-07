package errors

import "net/http"

type Type string

const (
	BAD_REQUEST      Type = "bad_request"
	INTERNAL_ERROR   Type = "internal_error"
	NOT_FOUND        Type = "not_found"
	VALIDATION_ERROR Type = "validation_error"
)

func (t Type) GetStatusCode() int {
	switch t {
	case BAD_REQUEST:
		return http.StatusBadRequest
	case INTERNAL_ERROR:
		return http.StatusInternalServerError
	case NOT_FOUND:
		return http.StatusNotFound
	case VALIDATION_ERROR:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
