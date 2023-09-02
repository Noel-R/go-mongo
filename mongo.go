package my_mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Context

type MongoContext struct {
	client        *mongo.Client
	clientOptions *options.ClientOptions
	context       context.Context
}

type IMongoContext interface {
	Create(string) *MongoContext
}

func (ctx *MongoContext) setOptions(connStr string) *options.ClientOptions {
	settings := options.Client()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	settings.SetServerAPIOptions(serverAPI)
	settings.ApplyURI(connStr)

	return settings
}

func (ctx *MongoContext) Create(connStr string) (*MongoContext, error) {
	var err error

	c := &MongoContext{}
	c.context = context.Background()
	c.clientOptions = c.setOptions(connStr)
	c.client, err = mongo.Connect(c.context, c.clientOptions)
	if err != nil {
		return nil, err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			panic(err)
		}
	}(c.client, c.context)

	if err = c.client.Database("talk-r").RunCommand(c.context, bson.D{{"ping", 1}}).Err(); err != nil {
		return nil, err
	}

	fmt.Println("Pinged db deployment.")
	return c, nil
}
