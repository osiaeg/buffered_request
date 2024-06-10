package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/osiaeg/buffered_request/internal/config"
	"github.com/osiaeg/buffered_request/internal/services"
	"github.com/osiaeg/buffered_request/internal/transport/rest"
)

func Run() {
	cfg := config.Parse("local")

	kafkaURL := fmt.Sprintf("%s:%s", cfg.Kafka.Host, cfg.Kafka.Port)

	p := services.NewProducer(kafkaURL, cfg.Kafka.Topic)
	mux := rest.Mux(p)

	fmt.Printf("Server launched.\nURL: http://%s:%s\n", cfg.Server.Host, cfg.Server.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Server.Port), mux))
}
