package json

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

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
