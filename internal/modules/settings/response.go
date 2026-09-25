package settings

import (
	"time"

	"github.com/google/uuid"
)

type SettingResponse struct {
	Key         string     `json:"key"`
	Value       string     `json:"value"`
	Type        string     `json:"type"`
	Description string     `json:"description"`
	Min         *float64   `json:"min,omitempty"`
	Max         *float64   `json:"max,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
	UpdatedBy   *uuid.UUID `json:"updatedBy,omitempty"`
}
