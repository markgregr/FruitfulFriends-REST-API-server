package producer

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"time"
)

type Producer struct {
	producer *kafka.Writer
	logger   *logrus.Logger
}

func NewProducerWorker(brokerAddr, topic string, logger *logrus.Logger) *Producer {
	return &Producer{
		producer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:  []string{brokerAddr},
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		}),
		logger: logger,
	}
}

func (pw *Producer) Start(ctx context.Context) error {
	const op = "producer.Producer.Start"
	pw.logger.WithField("operation", op).Info("Starting Kafka Producer Worker")

	for {
		select {
		case <-ctx.Done():
			pw.logger.WithField("operation", op).Info("Context done. Stopping Kafka Producer Worker.")
			return nil
		default:
			if err := pw.SendMessage("example message"); err != nil {
				pw.logger.WithError(err).Error("Failed to send message to Kafka")
			}
			time.Sleep(5 * time.Second)
		}
	}
}

func (pw *Producer) Stop() {
	// Метод остановки воркера
	// Здесь можно добавить код для остановки производства сообщений, если это необходимо
	// Пример: pw.logger.Info("Stopping Kafka Producer Worker")
	// pw.producer.Close()
}
