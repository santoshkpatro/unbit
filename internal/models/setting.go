package models

import (
	"encoding/json"
	"time"
)

type Setting struct {
	Key       string          `json:"key" db:"key"`
	Value     json.RawMessage `json:"value" db:"value"`
	UpdatedAt time.Time       `json:"updatedAt" db:"updated_at"`
}
