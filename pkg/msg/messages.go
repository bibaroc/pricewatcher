package msg

import (
	"context"
	"time"
)

type (
	Message struct {
		CreatedAt time.Time `json:"created_at,omitempty"`
		Text      string    `json:"text,omitempty"`
	}
	MessageRepository interface {
		Save(context.Context, string, string) error
		Get(context.Context, string) (*Message, error)
	}
)
