package kafka

import "github.com/Shopify/sarama"

type Producer interface {
	//PublishMessage(ctx context.Context, msgs ...kafka.Message) error
	PublishMessage(msgs ...*sarama.ProducerMessage) error
	Close() error
}
