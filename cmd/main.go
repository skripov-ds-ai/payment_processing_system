package main

import (
	"context"
	"errors"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

func main() {
	// make a writer that produces to topic-A, using the least-bytes distribution
	w := &kafka.Writer{
		Addr:                   kafka.TCP("localhost:9092", "localhost:9093", "localhost:9094"),
		Topic:                  "topic-A",
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	//err := w.WriteMessages(context.Background(),
	//	kafka.Message{
	//		Key:   []byte("Key-A"),
	//		Value: []byte("Hello World!"),
	//	},
	//	kafka.Message{
	//		Key:   []byte("Key-B"),
	//		Value: []byte("One!"),
	//	},
	//	kafka.Message{
	//		Key:   []byte("Key-C"),
	//		Value: []byte("Two!"),
	//	},
	//)
	//if err != nil {
	//	log.Fatal("failed to write messages:", err)
	//}
	//
	//if err := w.Close(); err != nil {
	//	log.Fatal("failed to close writer:", err)
	//}

	messages := []kafka.Message{
		{
			Key:   []byte("Key-A"),
			Value: []byte("Hello World!"),
		},
		{
			Key:   []byte("Key-B"),
			Value: []byte("One!"),
		},
		{
			Key:   []byte("Key-C"),
			Value: []byte("Two!"),
		},
	}
	var err error
	const retries = 3
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// attempt to create topic prior to publishing the message
		err = w.WriteMessages(ctx, messages...)
		if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
			time.Sleep(time.Millisecond * 250)
			continue
		}

		if err != nil {
			log.Fatalf("unexpected error %v", err)
		}
	}

	if err := w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
