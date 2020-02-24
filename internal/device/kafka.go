package device

import (
	"context"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog/log"
)

type kafkaConsumer interface {
	Poll(int) kafka.Event
	StoreOffsets(offsets []kafka.TopicPartition) (storedOffsets []kafka.TopicPartition, err error)
}

// KafkaSubscriber controls device information subscriptions
type KafkaSubscriber struct {
	kc       kafkaConsumer
	kmCh     chan *kafka.Message
	poolSize int
}

// NewKafkaSubscriber creates a new KafkaSubscriber
func NewKafkaSubscriber(kc kafkaConsumer) *KafkaSubscriber {
	return &KafkaSubscriber{
		kc:       kc,
		poolSize: 1000,
		kmCh:     make(chan *kafka.Message),
	}
}

// Start initialize a blocking subscriber
func (ks *KafkaSubscriber) Start(ctx context.Context) {
	log.Info().Msg("kafka device subscriber is listening for events...")

RUNNING:
	for {
		select {
		case <-ctx.Done():
			break RUNNING
		default:
			ev := ks.kc.Poll(ks.poolSize)
			switch e := ev.(type) {
			case kafka.Error:
				log.Error().Msgf("kafka error msg received: %s", e.Error())
				time.Sleep(time.Second)
			case *kafka.Message:
				ks.kmCh <- e
			}
		}
	}
}

// subscribe to consume received messages
func (ks *KafkaSubscriber) subscribe() <-chan *kafka.Message {
	return ks.kmCh
}

// commit acknowledge that the message received was correctly consumed
func (ks *KafkaSubscriber) commit(m *kafka.Message) error {
	if _, err := ks.kc.StoreOffsets([]kafka.TopicPartition{m.TopicPartition}); err != nil {
		return err
	}
	return nil
}
