package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func DecodeJSON(r *http.Request, dest interface{}) error {
	if r.Body == nil {
		return errors.New("request body is empty")
	}
	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		return fmt.Errorf("unexpected content type: %s", ct)
	}
	return json.NewDecoder(r.Body).Decode(dest)
}
