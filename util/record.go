package util

import (
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

//URL mamad
type URL struct {
	// ID primitive.ObjectID
	Protocol string
	Orig string
	Short string
	ExpTime time.Time
}

type ClientRequest struct {
	Url string
	Method string
} 