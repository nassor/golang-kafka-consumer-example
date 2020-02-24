package device

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockKafkaConsumer struct {
	mock.Mock
}

func (m *MockKafkaConsumer) Poll(s int) kafka.Event {
	args := m.Called(s)
	return args.Get(0).(kafka.Event)
}

func (m *MockKafkaConsumer) StoreOffsets(offsets []kafka.TopicPartition) (storedOffsets []kafka.TopicPartition, err error) {
	args := m.Called(offsets)
	return args.Get(0).([]kafka.TopicPartition), args.Error(1)
}

func TestNewKafkaSubscriber(t *testing.T) {
	ks := NewKafkaSubscriber(&MockKafkaConsumer{})
	assert.NotNil(t, ks.kc)
	assert.NotNil(t, ks.kmCh)
	assert.Equal(t, 1000, ks.poolSize)
}

func TestKafkaSubscriber_Start(t *testing.T) {
	testTable := []struct {
		name  string
		event kafka.Event
	}{
		{
			name:  "receive kafka message correctly",
			event: &kafka.Message{Key: []byte("ok")},
		},
		{
			name: "can receive a error event",
			// and will pause for one second before
			// poll the next messages.
			event: kafka.Error{},
		},
		{
			name: "start stop the context is cancelled",
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			mkc := &MockKafkaConsumer{}
			mkc.On("Poll", mock.AnythingOfType("int")).Return(tt.event)

			ks := NewKafkaSubscriber(mkc)
			ks.poolSize = 1
			go ks.Start(ctx)

			if _, ok := tt.event.(*kafka.Message); ok {
				msg := <-ks.kmCh
				assert.EqualValues(t, tt.event, msg)
			}
			if _, ok := tt.event.(kafka.Error); ok {
				time.Sleep(time.Millisecond)
			}

			cancel()
		})
	}
}

func TestKafkaSubscriber_subscribe(t *testing.T) {
	ks := NewKafkaSubscriber(&MockKafkaConsumer{})
	ch := ks.subscribe()
	assert.IsType(t, (<-chan *kafka.Message)(nil), ch)
}

func TestKafkaSubscriber_commit(t *testing.T) {
	testTable := []struct {
		name        string
		commitError error
	}{
		{
			name: "commit message correctly",
		},
		{
			name:        "commit can return an error",
			commitError: errors.New("an error"),
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			mkc := &MockKafkaConsumer{}
			msg := kafka.Message{TopicPartition: kafka.TopicPartition{}}
			mkc.On("StoreOffsets", []kafka.TopicPartition{msg.TopicPartition}).
				Return([]kafka.TopicPartition{}, tt.commitError)
			ks := NewKafkaSubscriber(mkc)
			err := ks.commit(&msg)
			assert.Equal(t, tt.commitError, err)
		})
	}
}
