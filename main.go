package main

import (
	"context"
	"fmt"
	"github.com/FerretDB/FerretDB/ferretdb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func main() {
	f, err := ferretdb.New(&ferretdb.Config{
		Listener: ferretdb.ListenerConfig{
			TCP: "127.0.0.1:27018",
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
		log.Fatal(err)
	}
	err = c.Ping(context.TODO(), nil)
	fmt.Println("sssss", err)

	cancel()
	<-done
}
