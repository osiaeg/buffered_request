package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/osiaeg/buffered_request/internal/config"
	"github.com/osiaeg/buffered_request/internal/services"
)

func Run() {
	cfg := config.Parse("local")

	kafkaURL := fmt.Sprintf("%s:%s", cfg.Kafka.Host, cfg.Kafka.Port)
	kafkaWriter := services.GetKafkaWriter(kafkaURL, cfg.Kafka.Topic)
	defer kafkaWriter.Close()

	kafkaReader := services.GetKafkaReader(
		kafkaURL,
		cfg.Kafka.Topic,
		cfg.Kafka.Partition,
	)
	defer kafkaReader.Close()

	go services.Consumer(kafkaReader)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /", services.ProducerHandler(kafkaWriter))

	fmt.Printf("Server launched.\nURL: http://%s:%s\n", cfg.Server.Host, cfg.Server.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Server.Port), mux))
}
