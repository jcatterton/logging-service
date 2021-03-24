package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"logging-service/pkg/api"
	"logging-service/pkg/consumer"
	"logging-service/pkg/dao"
	"logging-service/pkg/runner"
	"os"
)

func main() {
	broker := os.Getenv("BROKER")
	if broker == "" {
		logrus.Fatal("Broker is required")
		return
	}

	groupId := os.Getenv("GROUP_ID")
	if groupId == "" {
		logrus.Fatal("GroupID is required")
		return
	}

	topic := os.Getenv("TOPIC")
	if topic == "" {
		logrus.Fatal("Topic is required")
		return
	}

	c, err := consumer.CreateConsumer(broker, groupId, topic)
	if err != nil {
		logrus.WithError(err).Fatal("Error creating consumer")
		return
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		logrus.WithError(err).Fatal("Error creating mongo client")
		return
	}

	handler := dao.DbHandler{
		Client:     client,
		Database:   os.Getenv("DATABASE"),
		Collection: os.Getenv("COLLECTION"),
	}

	go runner.Run(c, &handler)

	if err := api.ListenAndServe(&handler); err != nil {
		logrus.WithError(err).Fatal("Error serving API")
	}
}
