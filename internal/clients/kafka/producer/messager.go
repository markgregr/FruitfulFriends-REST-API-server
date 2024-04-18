package producer

import (
	"context"
	"github.com/segmentio/kafka-go"
)

func (pw *Producer) SendMessage(message string) error {
	const op = "kafka.ProducerWorker.SendMessage"
	log := pw.logger.WithField("operation", op)
	log.Info(op, "Sending message to Kafka")

	err := pw.producer.WriteMessages(context.Background(), kafka.Message{Value: []byte(message)})
	if err != nil {
		pw.logger.WithError(err).Error("Failed to send message to Kafka")
		return err
	}

	return nil
}
