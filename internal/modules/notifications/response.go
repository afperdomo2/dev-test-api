package notifications

import (
	"time"

	"github.com/google/uuid"
)

type NotificationResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"userId"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	IsRead    bool      `json:"isRead"`
	CreatedAt time.Time `json:"createdAt"`
}

type UnreadCountResponse struct {
	Count int64 `json:"count"`
}

type MarkAllReadResponse struct {
	Updated int64 `json:"updated"`
}
