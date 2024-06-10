package services

import (
	"context"
	"errors"
	"log"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

func NewProducer(address, topic string) *producer {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(address),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
		Async:                  true,
	}

	return &producer{w}
}

type producer struct {
	w *kafka.Writer
}

func (p *producer) Produce(ctx context.Context, key, value []byte) error {
	msg := createKafkaMessage(key, value)
	const retries = 3

	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		err := p.w.WriteMessages(ctx, msg)
		if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
			time.Sleep(time.Millisecond * 250)
			log.Println(err)
			continue
		}

		if err != nil {
			log.Println("unexpected error %v", err)
			return err
		}
		break
	}

	return nil
}

func createKafkaMessage(key, value []byte) kafka.Message {
	return kafka.Message{
		Key:   key,
		Value: value,
	}
}
