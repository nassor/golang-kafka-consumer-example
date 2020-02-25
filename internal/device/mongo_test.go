package device

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MockMongoCollection struct {
	mock.Mock
}

func (m *MockMongoCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, document, opts)
	return nil, args.Error(0)
}

func (m *MockMongoCollection) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	args := m.Called(ctx, filter, opts)
	return nil, args.Error(0)
}

func TestMongoStore_add(t *testing.T) {
	testTable := []struct {
		name string
		err  error
	}{
		{
			name: "add device to the database",
		},
		{
			name: "insert can return an error",
			err:  errors.New("any"),
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mcoll := &MockMongoCollection{}
			mcoll.On("InsertOne",
				ctx,
				mock.AnythingOfType("Device"),
				mock.AnythingOfType("[]*options.InsertOneOptions")).
				Return(tt.err)

			mst := &MongoStore{
				c: mcoll,
			}
			mst.add(ctx, Device{ID: "123"}, time.Now())
			assert.True(t, mcoll.AssertExpectations(t))
		})
	}
}
