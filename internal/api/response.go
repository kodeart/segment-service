package api

import (
    "encoding/json"
    "errors"
    "net/http"

    "segment-service/internal/entities"
)

// jsonResponse writes a JSON response.
// The HTTP status code is derived from the error (if any),
// or overridden explicitly with the "status" argument.
//
// Usage:
//   jsonResponse(w, payload, nil)                     -> 200
//   jsonResponse(w, payload, nil, http.StatusCreated) -> 201
//   jsonResponse(w, nil, err)                         -> status derived from err
func jsonResponse(w http.ResponseWriter, payload any, err error, status ...int) {
    s := http.StatusOK
    if len(status) > 0 {
        s = status[0]
    }
    if err != nil {
        s = statusFromError(err)
        // RFC-7807
        data, _ := json.Marshal(map[string]any{"detail": err.Error(), "status": s})
        jsonWrite(w, s, data)
        return
    }
    data, _ := json.Marshal(payload)
    jsonWrite(w, s, data)
}

func jsonWrite(w http.ResponseWriter, status int, data []byte) {
    contentType := "application/json"
    if status >= http.StatusBadRequest {
        contentType = "application/problem+json" // RFC-7807
    }
    w.Header().Set("Content-Type", contentType)
    w.WriteHeader(status)
    // No response body for 204/304
    if status == http.StatusNoContent || status == http.StatusNotModified {
        return
    }
    if len(data) == 0 {
        return
    }
    _, _ = w.Write(data)
}

// statusFromError maps domain error to HTTP status code.
func statusFromError(err error) int {
    switch {
    case err == nil:
        return http.StatusOK
    case errors.Is(err, entities.ErrSegmentNotFound):
        return http.StatusNotFound
    case errors.Is(err, entities.ErrInvalidSegmentID),
        errors.Is(err, entities.ErrInvalidJsonBody),
        errors.Is(err, entities.ErrInvalidSegment):
        return http.StatusBadRequest
    default:
        return http.StatusUnprocessableEntity
    }
}
