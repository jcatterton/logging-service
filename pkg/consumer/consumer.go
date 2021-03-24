package consumer

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func CreateConsumer(broker string, groupId string, topic string) (*kafka.Consumer, error) {
	configMap := kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id":          groupId,
	}

	c, err := kafka.NewConsumer(&configMap)
	if err != nil {
		return nil, err
	}

	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return nil, err
	}

	return c, nil
}
