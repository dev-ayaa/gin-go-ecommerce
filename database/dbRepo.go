package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client = DatabaseSetUp()

func DatabaseSetUp() *mongo.Client {

	//setup the context
	ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelCtx()

	//creating new client
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)

	}
	//connecting the client to the database
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
		return client
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Println("Unable to ping the database")
		return nil
	}
	fmt.Println("Database Setup Successfully connected to Mongodb")
	return client

}

//UserData
func UserData(client *mongo.Client, collection string) *mongo.Collection {
	return nil
}

func ProductData(client *mongo.Client, collection string) *mongo.Collection {
	return nil
}
