package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
	"logging-service/pkg/dao"
	"logging-service/pkg/models"
	"os"
)

func Run(c *kafka.Consumer, handler *dao.DbHandler) {
	sigChan := make(chan os.Signal, 1)
	run := true
	logrus.Info("Starting consumer")

	for run {
		select {
		case sig := <-sigChan:
			logrus.WithField("signal", sig).Info("Terminating")
			run = false
		default:
			ev := c.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				logrus.WithFields(logrus.Fields{
					"message":   string(e.Value),
					"partition": e.TopicPartition,
					"offset":    e.TopicPartition.Offset,
				}).Debug("Message received")

				var log models.Log
				if err := json.NewDecoder(bytes.NewBuffer(e.Value)).Decode(&log); err != nil {
					logrus.WithError(err).Error("Error decoding message")
					continue
				}

				result, err := handler.AddLog(context.Background(), log)
				if err != nil {
					logrus.WithError(err).Error("Error adding log to database")
					continue
				}

				logrus.Infof("Log with ID '%v' added to database", result)
			case *kafka.Error:
				logrus.WithError(e).Error("Error received from kafka")
				if e.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
			default:
				logrus.WithField("event", e).Debug("Ignored event")
			}
		}
	}

	logrus.Info("Closing consumer")
	if err := c.Close(); err != nil {
		logrus.WithError(err).Error("Error closing consumer")
	}
}
