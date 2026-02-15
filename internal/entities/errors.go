package entities

import (
    "errors"
)

var (
    ErrSegmentNotFound  = errors.New("segment not found")
    ErrInvalidSegmentID = errors.New("invalid segment id")
    ErrInvalidJsonBody  = errors.New("invalid json body")
    ErrInvalidSegment   = errors.New("invalid segment")
)
