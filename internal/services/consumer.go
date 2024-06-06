package services

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"strings"

	kafka "github.com/segmentio/kafka-go"
)

func GetKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}

func Consumer(reader *kafka.Reader) {
	sender := NewSender()
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		res := &Request{}
		value := bytes.NewBuffer(m.Value)
		derr := json.NewDecoder(value).Decode(res)
		if derr != nil {
			panic(derr)
		}
		sender.SendRequest(res)
		// fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}
}
