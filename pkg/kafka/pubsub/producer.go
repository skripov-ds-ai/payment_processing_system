package pubsub

import (
	"github.com/Shopify/sarama"
	"go.uber.org/multierr"
)

//type producer struct {
//	log     *logger.Logger
//	address []string
//	w       *kafka.Writer
//	topic   string
//}
//
//type topicPartition struct {
//	topic     string
//	partition int32
//}

type publisher struct {
	kafkaVersion string
	appName      string
	address      []string
	producer     sarama.AsyncProducer
}

func NewPublisher(address []string, kafkaVersion, appName string) (*publisher, error) {
	//version, err := sarama.ParseKafkaVersion(kafkaVersion)
	//if err != nil {
	//	return nil, err
	//}
	producerConfig := sarama.NewConfig()
	producerConfig.Net.MaxOpenRequests = 1
	//producerConfig.Version = version
	producerConfig.Version = sarama.V2_8_1_0
	producerConfig.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	// TODO: uncomment these lines & handle successes and errors
	//producerConfig.Producer.Return.Errors = true
	//producerConfig.Producer.Return.Successes = true
	producerConfig.Producer.RequiredAcks = sarama.WaitForAll
	producerConfig.Producer.Idempotent = true
	producerConfig.Producer.Transaction.ID = appName
	p, err := sarama.NewAsyncProducer(address, producerConfig)
	if err != nil {
		return nil, err
	}
	return &publisher{
		appName: appName, kafkaVersion: kafkaVersion,
		address: address, producer: p,
	}, nil
}

func (p publisher) PublishMessage(msg ...*sarama.ProducerMessage) error {
	if err := p.producer.BeginTxn(); err != nil {
		return err
	}
	for _, m := range msg {
		p.producer.Input() <- m
	}
	if err := p.producer.CommitTxn(); err != nil {
		if err1 := p.producer.AbortTxn(); err1 != nil {
			multierr.AppendInto(&err, err1)
		}
		return err
	}
	return nil
}

func (p *publisher) Close() error {
	return p.producer.Close()
}

//// NewProducer create new kafka producer
//func NewProducer(address []string, topic string, log *logger.Logger) *producer {
//	writer := kafka.Writer{
//		Addr:        kafka.TCP(address...),
//		Topic:       topic,
//		Logger:      kafka.LoggerFunc(log.Infof),
//		ErrorLogger: kafka.LoggerFunc(log.Errorf),
//	}
//	return &producer{log: log, address: address, w: &writer, topic: topic}
//}
//
//func (p *producer) PublishMessage(ctx context.Context, msgs ...kafka.Message) error {
//	return p.w.WriteMessages(ctx, msgs...)
//}
//
//func (p *producer) Close() error {
//	return p.w.Close()
//}
