package api

import (
    "net/http"
)

// NewRouter creates a new HTTP router and registers
// segment-related routes with their corresponding handlers.
func NewRouter(hnd *SegmentHandler) http.Handler {
    mux := http.NewServeMux()

    mux.HandleFunc("GET /segments", hnd.handleList)
    mux.HandleFunc("GET /segments/{id}", hnd.handleGet)
    mux.HandleFunc("POST /segment", hnd.handleCreate)
    mux.HandleFunc("PUT /segments/{id}", hnd.handleUpdate)
    mux.HandleFunc("DELETE /segments/{id}", hnd.handleDelete)

    return mux
}
