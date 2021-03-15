package main

import (

	"context"
	// "encoding/json"
	"fmt"
	"log"
	// "net/http"
	"time"

	// "github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// TODO - change to string formatted values.
var (

	URI = "mongodb+srv://<username>:<password>@cicdcluster.4yy64.mongodb.net/<database>?retryWrites=true&w=majority"
)


type ListingAndReviews struct {

	ID 			primitive.ObjectID	`json:"_id,omitempty"`
	listing_url	string 				`json:"listing_url,omitempty"`
	name		string 				`json:"name,omitempty"`
	summary		string 				`json:"summary,omitempty"`
	description	string 				`json:"description,omitempty"`
}


func main() {

	fmt.Println("Connecting to Database..")

	// Setting the context and stacking it to close after main thread exits.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to the database. If there is an error,
	// check for error("err") and print.
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI))
	if err != nil {
		log.Println(err)
	}

	// Make sure to close the connection after the main thread is closed.
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// (Optional) Ping the primary - Checks if the mongo client is actually connected or not.
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected and pinged.")


	// get the collection
	collection := client.Database("sample_airbnb").Collection("listingsAndReviews")

	// Now process with the collection
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	count := 2

	for cursor.Next(ctx) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(result["_id"], result["amenities"])
		count--

		if count == 0 {
			break
		}
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
}
