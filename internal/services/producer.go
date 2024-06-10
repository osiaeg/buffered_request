package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

func checkBracket(decoder *json.Decoder) {
	_, err := decoder.Token()
	if err != nil {
		log.Fatalln(err)
	}
}

func createKafkaMessage(key, value []byte) kafka.Message {
	return kafka.Message{
		Key:   key,
		Value: value,
	}
}

func ProducerHandler(kafkaWriter *kafka.Writer) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(wrt http.ResponseWriter, req *http.Request) {
		log.Println("POST /")

		decoder := json.NewDecoder(req.Body)

		checkBracket(decoder)
		for decoder.More() {
			var request Request

			err := decoder.Decode(&request)
			if err != nil {
				http.Error(wrt, "Invalid data", http.StatusBadRequest)
				log.Println(err)
				continue
			}

			value, err := json.Marshal(request)
			if err != nil {
				log.Println(err)
			}

			key := []byte(fmt.Sprintf("address-%s", req.RemoteAddr))
			msg := createKafkaMessage(key, value)

			const retries = 3
			for i := 0; i < retries; i++ {
				log.Println(fmt.Sprintf("Try to send msg: %s", i))
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				err = kafkaWriter.WriteMessages(ctx, msg)
				if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
					time.Sleep(time.Millisecond * 250)
					continue
				}

				if err != nil {
					log.Fatalf("unexpected error %v", err)
				}
				break
			}

			// err := kafkaWriter.WriteMessages(req.Context(), msg)
			// if err != nil {
			// 	log.Fatalln(err)
			// }

			log.Println("Body from request succesfully written to kafka.")
		}
		checkBracket(decoder)
	})
}

func GetKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:                   kafka.TCP(kafkaURL),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
		Async:                  true,
	}
}
