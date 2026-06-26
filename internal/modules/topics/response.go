package topics

import (
	"time"

	"github.com/google/uuid"
)

type TopicResponse struct {
	ID        uuid.UUID  `json:"id"`
	Slug      string     `json:"slug"`
	Name      string     `json:"name"`
	Category  string     `json:"category"`
	IsSystem  bool       `json:"isSystem"`
	CreatedBy *uuid.UUID `json:"createdBy,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

type TopicListResponse struct {
	ID        uuid.UUID  `json:"id"`
	Slug      string     `json:"slug"`
	Name      string     `json:"name"`
	Category  string     `json:"category"`
	IsSystem  bool       `json:"isSystem"`
	CreatedBy *uuid.UUID `json:"createdBy,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
}
