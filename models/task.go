package models

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	IsCompleted   bool      `json:"is_completed"`
	AttachmentURL *string   `json:"attachment_url"`
	CreatedAt     time.Time `json:"created_at"`
}
