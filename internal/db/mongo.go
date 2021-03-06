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

	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI("mongodb://localhost:27017").SetAuth(credential))
	if err != nil {
		panic(err)
	}

	return client
}

func (mh *MongoHandler) SaveWebhook(webhook object.Webhook) primitive.ObjectID {
	client := mh.getClient()
	collection := client.Database(mh.Config.Database.Name).Collection("webhook")

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

func (mh *MongoHandler) SaveAuth(auth *object.Authentication) primitive.ObjectID {
	client := mh.getClient()
	collection := client.Database(mh.Config.Database.Name).Collection("auth")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc, err := ToDoc(auth)
	if err != nil {
		panic(err)
	}

	res, err := collection.InsertOne(ctx, doc)
	if err != nil {
		panic(err)
	}

	return res.InsertedID.(primitive.ObjectID)
}

func (mh *MongoHandler) GetAuth() *object.Authentication {
	client := mh.getClient()
	collection := client.Database(mh.Config.Database.Name).Collection("auth")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := collection.Find(
		ctx,
		bson.D{},
	)
	if err != nil {
		panic(err)
	}

	defer cur.Close(context.Background())

	var auth object.Authentication
	for cur.Next(context.Background()) {
		err = cur.Decode(&auth)
		if err != nil {
			panic(err)
		}
	}
	return &auth
}

func (mh *MongoHandler) DeleteAuth() error {
	client := mh.getClient()
	collection := client.Database(mh.Config.Database.Name).Collection("auth")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.DeleteMany(
		ctx,
		bson.D{},
	)
	if err != nil {
		return err
	}

	return nil
}

func (mh *MongoHandler) GetPendingWebhook() *object.Webhook {
	client := mh.getClient()
	collection := client.Database(mh.Config.Database.Name).Collection("webhook")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := collection.Find(
		ctx,
		bson.D{},
	)
	if err != nil {
		panic(err)
	}

	defer cur.Close(context.Background())

	var webhook object.Webhook
	for cur.Next(context.Background()) {
		err = cur.Decode(&webhook)
		if err != nil {
			panic(err)
		}
	}
	return &webhook
}

func ToDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)

	return
}
