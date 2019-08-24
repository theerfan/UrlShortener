package mongodb

//TODO: Change it to an api working with database, separated from the main server.

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/theerfan/urlshortener/util"
	"fmt"
	"time"
)

var collection *mongo.Collection


//GiveCount MAMAd
func GiveCount() int64 {
	ctx, _     := context.WithTimeout(context.Background(), 10*time.Second)
	count, err := collection.EstimatedDocumentCount(ctx)
	if err != nil {
		fmt.Println(err)
	}
	return count
}

//PutIntoDatabase mamad
func PutIntoDatabase(url util.URL) {
	fmt.Println(url.Protocol)
	if url.Protocol == "" {
		url.Protocol = "http"
	}
	ctx, _   := context.WithTimeout(context.Background(), 10*time.Second)
	bsonURL, err := bson.Marshal(url)
	collection.InsertOne(ctx, bsonURL)
	if err != nil {
		fmt.Println(err)
	}
}

func GetFromDatabase(short string) *util.URL {
	var full util.URL
	filter := bson.M{"short": short}
	ctx, _     := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, filter).Decode(&full)
	if err != nil {
		fmt.Println(err)
	}
	return &full
}

//Init initializes the database
func Init() {
	fmt.Println("mamaD")
	ctx, _      := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Println(err)
	}
	collection  = client.Database("url").Collection("short")
	fmt.Println("success")

}
