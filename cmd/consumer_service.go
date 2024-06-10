package main

import (
	"fmt"

	"github.com/osiaeg/buffered_request/internal/config"
	"github.com/osiaeg/buffered_request/internal/services"
)

func main() {
	cfg := config.Parse("local")

	kafkaURL := fmt.Sprintf("%s:%s", cfg.Kafka.Host, cfg.Kafka.Port)
	kafkaReader := services.GetKafkaReader(
		kafkaURL,
		cfg.Kafka.Topic,
	)
	defer kafkaReader.Close()

	services.Consumer(kafkaReader)
}
