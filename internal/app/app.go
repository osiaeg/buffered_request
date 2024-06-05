package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/osiaeg/buffered_request/internal/config"
	"github.com/osiaeg/buffered_request/internal/services"
	"github.com/segmentio/kafka-go"
)

func saveAct(w http.ResponseWriter, r *http.Request) {
	// sender := services.NewSender()
	// sender.SendRequest()
	// services.RunProducer()
}

func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}

func Run() {
	cfg := config.Parse("local")

	kafkaWriter := services.GetKafkaWriter("localhost:9092", "task")
	defer kafkaWriter.Close()

	reader := getKafkaReader("localhost:9092", "task", "0")
	defer reader.Close()

	fmt.Println("start consuming ... !!")
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /save_acts", saveAct)
	mux.HandleFunc("POST /", services.ProducerHandler(kafkaWriter))

	fmt.Printf("Server launched.\nURL: http://%s:%s\n", cfg.Server.Host, cfg.Server.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Server.Port), mux))
}
