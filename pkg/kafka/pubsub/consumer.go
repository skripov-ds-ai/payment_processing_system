package pubsub

import "github.com/Shopify/sarama"

type reader struct {
	address       []string
	groupID       string
	consumerGroup sarama.ConsumerGroup
}

//func NewReader(address []string, groupID string) (*reader, error) {
func NewReader(address []string, groupID string) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	//config.ClientID = clientID
	// TODO: move from hardcode
	config.Version = sarama.V2_8_1_0
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.IsolationLevel = sarama.ReadCommitted

	// TODO: add client.Close() to external
	consumer, err := sarama.NewConsumerGroup(address, groupID, config)
	if err != nil {
		return nil, err
	}
	return consumer, nil
	//return &reader{
	//	consumerGroup: consumer,
	//	address:       address, groupID: groupID,
	//}, nil
}

//type consumer struct {
//	r       *kafka.Reader
//	log     *logger.Logger
//	address []string
//	topic   string
//}
//
//func NewConsumer(address []string, topic, groupID string, log *logger.Logger) *consumer {
//	// TODO: check NewConsumerGroup - https://github.com/segmentio/kafka-go/blob/294fbdbf43ce5c2bc6aa89d8a35db949e1f95367/example_consumergroup_test.go#L11
//	r := kafka.NewReader(kafka.ReaderConfig{
//		Brokers:        address,
//		Topic:          topic,
//		GroupID:        groupID,
//		CommitInterval: time.Second,
//		Logger:         kafka.LoggerFunc(log.Infof),
//		ErrorLogger:    kafka.LoggerFunc(log.Errorf),
//		MinBytes:       10e3, // 10KB
//		MaxBytes:       10e6, // 10MB
//	})
//	return &consumer{r: r, log: log, address: address, topic: topic}
//}
//
//func (c *consumer) ReadMessage(ctx context.Context) (kafka.Message, error) {
//	return c.r.ReadMessage(ctx)
//}
//
//func (c *consumer) Close() error {
//	return c.r.Close()
//}
