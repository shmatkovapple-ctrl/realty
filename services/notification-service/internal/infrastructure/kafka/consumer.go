package kafka

import (
"context"
"encoding/json"
"log"

"github.com/segmentio/kafka-go"
)

type Handler func(ctx context.Context, event map[string]any) error

type Consumer struct {
reader  *kafka.Reader
handler Handler
}

func NewConsumer(brokerURL, topic, groupID string, handler Handler) *Consumer {
return &Consumer{
reader: kafka.NewReader(kafka.ReaderConfig{
Brokers: []string{brokerURL},
Topic:   topic,
GroupID: groupID,
}),
handler: handler,
}
}

func (c *Consumer) Start(ctx context.Context) {
go func() {
for {
msg, err := c.reader.ReadMessage(ctx)
if err != nil {
if ctx.Err() != nil {
return
}
log.Printf("ошибка чтения Kafka: %v", err)
continue
}

var event map[string]any
if err := json.Unmarshal(msg.Value, &event); err != nil {
log.Printf("ошибка десериализации события: %v", err)
continue
}

if err := c.handler(ctx, event); err != nil {
log.Printf("ошибка обработки события: %v", err)
}
}
}()
}

func (c *Consumer) Close() error {
return c.reader.Close()
}
