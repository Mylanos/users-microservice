package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func WriteResponse(w http.ResponseWriter, status int, payload []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(payload)
}

func ConstructResponse(w http.ResponseWriter, status int, data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		log.Print("ERROR: failed to serialize the response")
		w.WriteHeader(http.StatusInternalServerError)
		return NewAPIError(http.StatusInternalServerError, "failed to serialize the response")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(payload); err != nil {
		log.Printf("ERROR: Failed to write payload to client: %v", err)
	}
	return nil
}

func ConstructResponseWithError(w http.ResponseWriter, err APIError) error {
	response := APIResponse{
		Success:   false,
		Data:      nil,
		Error:     &err,
		Timestamp: time.Now(),
	}
	return ConstructResponse(w, err.Status, response)
}

func ConstructSuccessResponse(w http.ResponseWriter, status int, data UserAPI) error {
	response := APIResponse{
		Success:   true,
		Data:      data,
		Error:     nil,
		Timestamp: time.Now(),
	}
	return ConstructResponse(w, status, response)
}

func logError(r *http.Request, err error, duration time.Duration) {
	log.Printf("ERROR: %s %s - %v (took %v)",
		r.Method,
		r.URL.Path,
		err,
		duration,
	)
}

func logSuccess(r *http.Request, duration time.Duration) {
	log.Printf("SUCCESS: %s %s (took %v)",
		r.Method,
		r.URL.Path,
		duration,
	)
}
