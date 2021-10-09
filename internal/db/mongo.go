package db

import (
	"context"
	"time"

	"github.com/m4tthewde/fdmxyz/internal/config"
	"github.com/m4tthewde/fdmxyz/internal/object"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoHandler struct {
	Config *config.Config
}

func (mh *MongoHandler) getClient() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	credential := options.Credential{
		Username: mh.Config.Database.User,
		Password: mh.Config.Database.Password,
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017").SetAuth(credential))
	if err != nil {
		panic(err)
	}

	return client
}

func (mh *MongoHandler) SaveWebhook(webhook object.Webhook) primitive.ObjectID {
	client := mh.getClient()
	collection := client.Database(mh.Config.Database.Name).Collection("endResult")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc, err := ToDoc(webhook)
	if err != nil {
		panic(err)
	}

	res, err := collection.InsertOne(ctx, doc)
	if err != nil {
		panic(err)
	}

	return res.InsertedID.(primitive.ObjectID)
}

func ToDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)

	return
}
