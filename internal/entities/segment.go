package entities

import (
    "time"
)

type Segment struct {
    ID         int        `json:"id" db:"id"`
    Name       string     `json:"name" db:"name"`
    TTLSeconds *int       `json:"ttl_seconds,omitempty" db:"ttl_seconds"`
    CreateAt   time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
    DeletedAt  *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}
