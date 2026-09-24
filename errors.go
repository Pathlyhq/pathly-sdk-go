package pathly

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// APIError carries the status and the message returned by the API.
//
// The API messages are written to be displayed as they are: rewording them
// would lose the detail that makes the fix possible (the named missing scope,
// the offending field, the exceeded limit).
type APIError struct {
	StatusCode int
	Message    string
	Path       string
}

func (e *APIError) Error() string {
	if e.Path == "" {
		return fmt.Sprintf("%s (HTTP %d)", e.Message, e.StatusCode)
	}
	return fmt.Sprintf("%s: %s (HTTP %d)", e.Path, e.Message, e.StatusCode)
}

// IsNotFound reports a resource missing from this organization.
func (e *APIError) IsNotFound() bool { return e.StatusCode == http.StatusNotFound }

// IsNotFound is true when the error reports a missing resource.
func IsNotFound(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.IsNotFound()
}

func apiError(status int, path string, body []byte) *APIError {
	var parsed struct {
		Error   string `json:"error"`
		Details []struct {
			Message string `json:"message"`
			Path    []any  `json:"path"`
		} `json:"details"`
	}
	message := strings.TrimSpace(string(body))
	if err := json.Unmarshal(body, &parsed); err == nil && parsed.Error != "" {
		message = parsed.Error
		for _, d := range parsed.Details {
			message += fmt.Sprintf(" [%s: %s]", fieldOf(d.Path), d.Message)
		}
	}
	if message == "" {
		message = http.StatusText(status)
	}
	switch status {
	case http.StatusUnauthorized:
		message += " Check PATHLY_API_TOKEN: an expired, revoked or truncated key gives the same response."
	case http.StatusForbidden:
		message += " Widen the scopes of the key, or check that the plan includes the feature."
	}
	return &APIError{StatusCode: status, Message: message, Path: path}
}

// fieldOf renders a validation path readable: `events.0`, or `body` when empty.
func fieldOf(path []any) string {
	parts := make([]string, 0, len(path))
	for _, p := range path {
		switch v := p.(type) {
		case string:
			parts = append(parts, v)
		case float64:
			parts = append(parts, strconv.Itoa(int(v)))
		}
	}
	if len(parts) == 0 {
		return "body"
	}
	return strings.Join(parts, ".")
}
