package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"logging-service/pkg/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Handler interface {
	Ping(ctx context.Context) error
	AddLog(ctx context.Context, log models.Log) (interface{}, error)
	GetLogs(ctx context.Context) ([]models.Log, error)
}

type DbHandler struct {
	Client     *mongo.Client
	Database   string
	Collection string
}

func (db *DbHandler) Ping(ctx context.Context) error {
	return db.Client.Ping(ctx, readpref.Primary())
}

func (db *DbHandler) AddLog(ctx context.Context, log models.Log) (interface{}, error) {
	result, err := db.getCollection().InsertOne(ctx, log)
	if err != nil {
		return "", err
	}

	results, err := db.GetLogs(ctx)
	if err != nil {
		return "", err
	}

	if len(results) > 500 {
		deleteOpts := options.FindOneAndDeleteOptions{}
		deleteOpts.SetSort(bson.D{{"timeStamp", 1}})

		result := db.getCollection().FindOneAndDelete(ctx, bson.D{}, &deleteOpts)
		if result.Err() != nil {
			return "", result.Err()
		}
	}

	return result.InsertedID, nil
}

func (db *DbHandler) GetLogs(ctx context.Context) ([]models.Log, error) {
	cur, err := db.getCollection().Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var results []models.Log
	if err = cur.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (db *DbHandler) getCollection() *mongo.Collection {
	return db.Client.Database(db.Database).Collection(db.Collection)
}
