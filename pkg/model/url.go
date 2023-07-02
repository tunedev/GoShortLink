package model

import (
	"time"
)

type Url struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	ShortUrl  string    `json:"short_url" bson:"short_url,omitempty"`
	LongUrl   string    `json:"long_url" bson:"long_url,omitempty"`
	CreatedAt time.Time `json:"created_at" bson:"created_at,omitempty"`
}