/*
Package "providers" is a package that provides required things reqired by the bot
to be used by other funcs.
This package should has all providers for the bot.
*/
package mongoProvider

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"GoAntispamBot/bot"
	"GoAntispamBot/bot/helpers/errHandler"
)

var ctx = context.Background()

func connect() (*mongo.Database, error) {
	// Initiate MongoDB connection.
	clientOptions := options.Client()
	clientOptions.ApplyURI(bot.BotConfig.DatabaseURL)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	// Connect to MongoDB.
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return client.Database("antispambot"), nil
}

func Insert(coll string, doc interface{}) {
	db, err := connect() // Init connection
	errHandler.Fatal(err)

	// Start inserting..
	_, err = db.Collection(coll).InsertOne(ctx, doc)
	errHandler.Fatal(err)
}

func Update(coll string, filter interface{}, update interface{}, upsert bool) {
	db, err := connect() // Init connection
	errHandler.Fatal(err)

	// Start updating...
	upserts := options.Update().SetUpsert(upsert)
	_, err = db.Collection(coll).UpdateOne(ctx, filter, update, upserts)
	errHandler.Fatal(err)
}

func Remove(coll string, doc interface{}) {
	db, err := connect() // Init connection
	errHandler.Fatal(err)

	// Start deleting...
	_, err = db.Collection(coll).DeleteOne(ctx, doc)
	errHandler.Fatal(err)
}

func FindOne(coll string, doc interface{}) (ret bson.Raw) {
	db, err := connect() // Init connection
	errHandler.Fatal(err)

	// Start searching...
	csr := db.Collection(coll).FindOne(ctx, doc)
	ret, err = csr.DecodeBytes()
	errHandler.Fatal(err)
	return
}
