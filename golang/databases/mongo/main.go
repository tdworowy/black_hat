package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Transaction struct {
	CCNum      string  `bson:"ccnum"`
	Date       string  `bson:"date"`
	Amount     float64 `bson:"amount"`
	Cvv        string  `bson:"cvv"`
	Expiration string  `bson:"exp"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	collection := client.Database("store").Collection("transactions")
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Panicln(err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var r Transaction
		err := cur.Decode(&r)
		if err != nil {
			log.Panicln(err)
		}
		fmt.Println(r)

	}

}
