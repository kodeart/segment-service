package api

import (
    "net/http"
    "strconv"

    "segment-service/internal/entities"
    "segment-service/internal/service"
)

type SegmentHandler struct {
    segmentService *service.SegmentsService
}

func NewSegmentHandler(segmentService *service.SegmentsService) *SegmentHandler {
    return &SegmentHandler{segmentService: segmentService}
}

// handleList all active segments.
// GET /segments
func (h *SegmentHandler) handleList(w http.ResponseWriter, r *http.Request) {
    list, err := h.segmentService.List(r.Context())
    jsonResponse(w, list, err)
}

// handleGet segment by ID.
// GET /segments/{id}
func (h *SegmentHandler) handleGet(w http.ResponseWriter, r *http.Request) {
    id, err := getIDFromPath(r)
    if err != nil {
        jsonResponse(w, nil, err)
        return
    }
    segment, err := h.segmentService.Get(r.Context(), id)
    jsonResponse(w, segment, err)
}

// handleCreate creates a new segment.
// POST /segment
func (h *SegmentHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
    var req segmentRequestInput
    if err := validateSegmentInput(r, &req); err != nil {
        jsonResponse(w, nil, err)
        return
    }

    segment := &entities.Segment{
        Name:       req.Name,
        TTLSeconds: req.TTLSeconds,
    }
    segment, err := h.segmentService.Create(r.Context(), segment)
    jsonResponse(w, segment, err, http.StatusCreated)
}

// handleUpdate updates a segment.
// PUT /segments/{id}
func (h *SegmentHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
    id, err := getIDFromPath(r)
    if err != nil {
        jsonResponse(w, nil, err)
        return
    }

    var req segmentRequestInput
    if err = validateSegmentInput(r, &req); err != nil {
        jsonResponse(w, nil, err)
        return
    }

    segment := &entities.Segment{
        ID:         id,
        Name:       req.Name,
        TTLSeconds: req.TTLSeconds,
    }
    updated, err := h.segmentService.Update(r.Context(), segment)
    jsonResponse(w, updated, err)
}

// handleDelete soft deletes a segment.
// DELETE /segments/{id}
func (h *SegmentHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
    id, err := getIDFromPath(r)
    if err != nil {
        jsonResponse(w, nil, err)
        return
    }
    err = h.segmentService.Delete(r.Context(), id)
    jsonResponse(w, nil, err, http.StatusNoContent)
}

// getIDFromPath is a helper method for
// parsing the segment ID from the URL,
func getIDFromPath(r *http.Request) (int, error) {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil {
        return 0, entities.ErrInvalidSegmentID
    }
    return id, nil
}
