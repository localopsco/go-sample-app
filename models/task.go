package models

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	IsCompleted   bool      `json:"bool"`
	AttachmentURL string    `json:"attachement_url"`
	CreatedAt     time.Time `json:"created_at"`
}
