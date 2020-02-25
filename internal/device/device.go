package device

import "time"

// Device is a mysterious hardware
type Device struct {
	ID        string    `bson:"id"`
	Version   string    `bson:"version"`
	Reference string    `bson:"ref"`
	Removed   bool      `bson:"-"`
	Disabled  bool      `bson:"-"`
	CreatedAt time.Time `bson:"created_at"`
}
