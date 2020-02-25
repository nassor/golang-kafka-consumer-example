package device

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoCollection interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
}

// MongoStore for device data
type MongoStore struct {
	c mongoCollection
}

// NewMongoStore for devices
func NewMongoStore(deviceCollection mongoCollection) *MongoStore {
	return &MongoStore{
		c: deviceCollection,
	}
}

func (st *MongoStore) add(ctx context.Context, d Device, now time.Time) error {
	if _, err := st.c.InsertOne(ctx, d); err != nil {
		return fmt.Errorf("adding device: %w", err)
	}
	// st.removeOld(ctx, d.ID, now)
	return nil
}

const expireInDays = 30

// func (st *MongoStore) removeOld(ctx context.Context, id string, now time.Time) error {
// 	expirationDate := now.AddDate(0, 0, (expireInDays * -1))
// 	filter := bson.D{
// 		{"id", id},
// 		{"created_at", bson.D{
// 			{"$lte", expirationDate},
// 		}},
// 	}
// 	if _, err := st.c.DeleteMany(ctx, filter); err != nil {
// 		return fmt.Errorf("removing devices older than %s: %w", expirationDate, err)
// 	}
// 	return nil
// }
