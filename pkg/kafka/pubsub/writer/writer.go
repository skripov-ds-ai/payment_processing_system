package writer

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type Writer struct {
	urls []string
	w    *kafka.Writer
}

func NewWriter(urls []string) *Writer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(urls...),
		RequiredAcks: kafka.RequireAll, // TODO: fix
		Balancer:     &kafka.LeastBytes{},
	}
	return &Writer{
		w:    w,
		urls: urls,
	}
}

func (w *Writer) Publish(ctx context.Context, msgs ...kafka.Message) error {
	return w.w.WriteMessages(ctx, msgs...)
}

func (w *Writer) Close() error {
	return w.w.Close()
}
