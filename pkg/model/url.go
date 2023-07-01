package model

import "time"

type Url struct {
	ID string `bson:"_id"`
	ShortUrl string `bson:"short_url"`
	LongUrl string `bson:"long_url"`
	CreatedAt time.Time `bson:"created_at"`
}