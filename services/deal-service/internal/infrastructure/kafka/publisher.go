package kafka

import (
"context"
"encoding/json"
"fmt"

"github.com/segmentio/kafka-go"
)

type Publisher struct {
writer *kafka.Writer
}

func NewPublisher(brokerURL string) *Publisher {
return &Publisher{
writer: &kafka.Writer{
Addr:     kafka.TCP(brokerURL),
Balancer: &kafka.LeastBytes{},
},
}
}

func (p *Publisher) Publish(ctx context.Context, topic string, event any) error {
data, err := json.Marshal(event)
if err != nil {
return fmt.Errorf("сериализация события: %w", err)
}
return p.writer.WriteMessages(ctx, kafka.Message{
Topic: topic,
Value: data,
})
}

func (p *Publisher) Close() error {
return p.writer.Close()
}
