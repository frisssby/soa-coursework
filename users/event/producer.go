package event

import (
	"encoding/json"

	"users/models"

	"github.com/IBM/sarama"
)

var eventProducer sarama.SyncProducer

func InitEventProducer(uri string) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	var err error
	eventProducer, err = sarama.NewSyncProducer([]string{uri}, config)
	return err
}

func ProduceEvent(topic string, event models.Event) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: 0,
		Value:     sarama.ByteEncoder(bytes),
	}
	_, _, err = eventProducer.SendMessage(msg)
	return err
}
