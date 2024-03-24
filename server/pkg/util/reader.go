package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func UnmarshalHTTPResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status %d but got %d: %s", http.StatusOK, resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("error unmarshalling response body: %w", err)
	}

	return nil
}

func JsonResponse(w http.ResponseWriter, code int, obj any) {
	b, err := json.Marshal(obj)
	if err != nil {
		log.Errorf("failed to send json response: %v", err)
		w.Write([]byte("Something went wrong"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	w.Write(b)

}

// logs the error and responses to caller
func ErrResponse(w http.ResponseWriter, code int, err error, msg string) {
	log.Error(err)
	w.WriteHeader(code)
	w.Write([]byte(msg))
}

func Response(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	w.Write([]byte(msg))
}
