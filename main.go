package main

import (
	"context"
	"fmt"
	"github.com/FerretDB/FerretDB/ferretdb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func insertData(client *mongo.Client) error {
	// Get a handle to the "test" database and the "example" collection
	collection := client.Database("test").Collection("example")

	// Create a document (a map) with the data you want to insert
	data := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}

	// Insert the document into the collection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	fmt.Println("Data inserted successfully!")
	return nil
}

func main() {
	f, err := ferretdb.New(&ferretdb.Config{
		Listener: ferretdb.ListenerConfig{
			TCP: "127.0.0.1:17028",
		},
		Handler:       "postgresql",
		PostgreSQLURL: "postgres://127.0.0.1:5432/ferretdb",
	})
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})

	go func() {
		log.Print(f.Run(ctx))
		close(done)
	}()

	uri := f.MongoDBURI()
	fmt.Println(uri)

	c, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Println("con", err)
	}
	err = c.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println("ping", err)
	}
	err = insertData(c)
	if err != nil {
		fmt.Println("insert", err)
	}

	// Specify the name of the database and collection
	databaseName := "test"      // Replace with your actual database name
	collectionName := "example" // Replace with the name of your collection

	fmt.Println(c.ListDatabases(ctx, bson.D{{}}))

	// Get a handle to the specified database and collection
	db := c.Database(databaseName)
	cl, e := db.ListCollectionNames(ctx, bson.D{{}})
	if e != nil {
		fmt.Println("c", e)
	}
	fmt.Println(cl)
	collection := db.Collection(collectionName)

	// Define a filter to find documents
	filter := bson.D{{}}
	// Replace "Field1" and "value_to_find" with the field and value you want to search for.

	// Find the documents that match the filter
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	// Iterate through the results
	results := []map[string]interface{}{}
	for cursor.Next(context.Background()) {
		var result map[string]interface{}
		err := cursor.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, result)
	}

	// Check for errors during the cursor iteration
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor
	cursor.Close(context.Background())

	// Print the found documents
	fmt.Println("Found documents:")
	fmt.Println(len(results))

	/*	db, err = c.ListDatabaseNames(context.TODO(), nil)
		if err != nil {
			fmt.Println("db", err)
		}
		fmt.Println(db)*/

	cancel()
	<-done
}
