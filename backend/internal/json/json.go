package json

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var Validate = validator.New()

// ValidUUID reports whether id is a well-formed UUID. When it is not, it writes
// a 400 response and returns false. Use it to guard handlers that take an {id}
// path parameter so a malformed id yields 400 instead of a generic 500.
func ValidUUID(w http.ResponseWriter, id string) bool {
	if err := uuid.Validate(id); err != nil {
		WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid id: %q", id))
		return false
	}
	return true
}

type ErrorResponse struct {
	Error     string    `json:"error"`
	Timestamp time.Time `json:"timestamp"`
}

func WriteError(w http.ResponseWriter, code int, error error) {
	WriteJSON(w, code, ErrorResponse{
		Error:     error.Error(),
		Timestamp: time.Now(),
	})
}

func WriteJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func Read(r *http.Request, data any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func ParseInt32(value string, fallback int32) int32 {
	var result int32
	_, err := fmt.Sscan(value, &result)
	if err != nil {
		return fallback
	}
	return result
}
